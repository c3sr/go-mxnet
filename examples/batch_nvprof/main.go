// +build !nogpu

package main

// #cgo linux CFLAGS: -I/usr/local/cuda/include
// #cgo linux LDFLAGS: -lcuda -lcudart -L/usr/local/cuda/lib64
// #include <cuda.h>
// #include <cuda_runtime.h>
// #include <cuda_profiler_api.h>
import "C"

import (
	"bufio"
	"context"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/k0kubun/pp"
	"github.com/c3sr/config"
	"github.com/c3sr/dlframework"
	"github.com/c3sr/dlframework/framework/feature"
	"github.com/c3sr/dlframework/framework/options"
	"github.com/c3sr/downloadmanager"
	"github.com/c3sr/go-mxnet/mxnet"
	nvidiasmi "github.com/c3sr/nvidia-smi"
  _ "github.com/c3sr/tracer/all"
	gotensor "gorgonia.org/tensor"
)

// https://github.com/dmlc/gluon-cv/blob/master/gluoncv/data/transforms/presets/imagenet.py
// mean=(0.485, 0.456, 0.406), std=(0.229, 0.224, 0.225)
var (
	batchSize   = 1
	model       = "squeezenet1.0"
	shape       = []int{1, 3, 224, 224}
	mean        = []float32{0.485, 0.456, 0.406}
  scale       = []float32{0.229, 0.224, 0.225}
	imgDir, _   = filepath.Abs("../_fixtures")
	imgPath     = filepath.Join(imgDir, "platypus.jpg")
	graph_url   = "http://s3.amazonaws.com/store.carml.org/models/mxnet/gluoncv/squeezenet1.0/model-symbol.json"
	weights_url = "http://s3.amazonaws.com/store.carml.org/models/mxnet/gluoncv/squeezenet1.0/model-0000.params"
	synset_url  = "http://s3.amazonaws.com/store.carml.org/synsets/imagenet/synset.txt"
)

// convert go Image to 1-dim array
func cvtRGBImageToNCHW1DArray(src image.Image, mean []float32, scale []float32) ([]float32, error) {
	if src == nil {
		return nil, fmt.Errorf("src image nil")
	}

	in := src.(*types.RGBImage)
	height := in.Bounds().Dy()
	width := in.Bounds().Dx()

	out := make([]float32, 3*height*width)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := y*in.Stride + x*3
			rgb := in.Pix[offset : offset+3]
			r, g, b := rgb[0], rgb[1], rgb[2]
			out[y*width+x] = (float32(r)/255 - mean[0]) / scale[0]
			out[width*height+y*width+x] = (float32(g)/255 - mean[1]) / scale[1]
			out[2*width*height+y*width+x] = (float32(b)/255 - mean[2]) / scale[2]
		}
	}

	return out, nil
}

func main() {
	dir, _ := filepath.Abs("../tmp")
	dir = filepath.Join(dir, model)
	graph := filepath.Join(dir, "model-symbol.json")
	weights := filepath.Join(dir, "model-0000.params")
	synset := filepath.Join(dir, "synset.txt")

	if !com.IsFile(graph) {
		if _, err := downloadmanager.DownloadInto(graph_url, dir); err != nil {
			panic(err)
		}
	}
	if !com.IsFile(weights) {
		if _, err := downloadmanager.DownloadInto(weights_url, dir); err != nil {
			panic(err)
		}
	}
	if !com.IsFile(synset) {
		if _, err := downloadmanager.DownloadInto(synset_url, dir); err != nil {
			panic(err)
		}
	}

	// load model
	symbol, err := ioutil.ReadFile(graph)
	if err != nil {
		panic(err)
	}
	params, err := ioutil.ReadFile(weights)
	if err != nil {
		panic(err)
	}


	height := shape[2]
	width := shape[3]
	channels := shape[1]

	r, err := os.Open(imgPath)
	if err != nil {
		panic(err)
	}

	var imgOpts []raiimage.Option
	imgOpts = append(imgOpts, raiimage.Mode(types.RGBMode))
	img, err := raiimage.Read(r, imgOpts...)
	if err != nil {
		panic(err)
	}

	imgOpts = append(imgOpts, raiimage.Resized(height, width))
	imgOpts = append(imgOpts, raiimage.ResizeAlgorithm(types.ResizeAlgorithmLinear))
	resized, err := raiimage.Resize(img, imgOpts...)
	if err != nil {
		panic(err)
	}

	input := make([]*gotensor.Dense, batchSize)
	imgFloats, err := cvtRGBImageToNCHW1DArray(resized, mean, scale)
	if err != nil {
		panic(err)
	}

	for ii := 0; ii < batchSize; ii++ {
		input[ii] = gotensor.New(
			gotensor.Of(tensor.Float32),
			gotensor.WithShape(height, width, channels),
			gotensor.WithBacking(imgFloats),
		)
	}

	device := options.CPU_DEVICE
	if nvidiasmi.HasGPU {
		device = options.CUDA_DEVICE

	}

	ctx := context.Background()

	in := options.Node{
		Key:   "data",
		Shape: shape,
	}
	predictor, err := mxnet.New(
		ctx,
		options.WithOptions(opts),
		options.Device(device, 0),
		options.Graph(symbol),
		options.Weights(params),
		options.BatchSize(batchSize),
		options.InputNodes([]options.Node{in}),
		options.OutputNodes([]options.Node{
			options.Node{Dtype: tensor.Float32},
		}),
	)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	defer predictor.Close()

  inputs := []tensor.Tensor{
		tensor.NewDense(tensor.Float32, in.Shape, tensor.WithBacking(input)),
  }
  
	C.cudaProfilerStart()

	err = predictor.Predict(ctx, inputs)
	if err != nil {
		panic(err)
	}

	C.cudaDeviceSynchronize()
	C.cudaProfilerStop()

	outputs, err := predictor.ReadPredictionOutputs(ctx)
	if err != nil {
		panic(err)
	}

  if len(outputs) != 1 {
		panic(errors.Errorf("invalid output length. got outputs of length %v", len(outputs)))
	}

  output := outputs[0].Data().([]float32)

	var labels []string
	f, err := os.Open(synset)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		labels = append(labels, line)
	}

	features := make([]dlframework.Features, batchSize)
	featuresLen := len(output) / batchSize

	for ii := 0; ii < batchSize; ii++ {
    rprobs := make([]*dlframework.Feature, featuresLen)
		soutputs := output[ii*featuresLen : (ii+1)*featuresLen]
		for jj := 0; jj < featuresLen; jj++ {
			rprobs[jj] = feature.New(
				feature.ClassificationIndex(int32(jj)),
				feature.ClassificationLabel(labels[jj]),
				feature.Probability(output[ii*featuresLen+jj]),
			)
    }
    nprobs := dlframework.Features(rprobs).ProbabilitiesApplySoftmaxFloat32()
		sort.Sort(nprobs)
		features[ii] = nprobs
	}

	if true {
		for i := 0; i < 1; i++ {
			results := features[i]
			top1 := results[0]
			pp.Println(top1.Probability)
			pp.Println(top1.GetClassification().GetLabel())
		}
	} else {
		_ = features
	}
}

func init() {
	config.Init(
		config.AppName("carml"),
		config.VerboseMode(true),
		config.DebugMode(true),
	)
}

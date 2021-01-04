# go-mxnet

[![Build Status](https://dev.azure.com/yhchang/c3sr/_apis/build/status/c3sr.go-mxnet?branchName=master)](https://dev.azure.com/yhchang/c3sr/_build/latest?definitionId=6&branchName=master)
[![Build Status](https://dev.azure.com/dakkak/rai/_apis/build/status/rai-project.go-mxnet)](https://dev.azure.com/dakkak/rai/_build/latest?definitionId=8)
[![Go Report Card](https://goreportcard.com/badge/github.com/c3sr/go-mxnet)](https://goreportcard.com/report/github.com/c3sr/go-mxnet)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Go binding for MXNet C predict API.
This is used by the [MXNet agent](https://github.com/c3sr/mxnet) in [MLModelScope](mlmodelscope.org) to perform model inference in Go.

## Installation

Download and install go-mxnet:

```
go get -v github.com/c3sr/go-mxnet
```

The binding requires MXNet and other Go packages.

### MXNet C Library

The MXNet C library is expected to be under `/opt/mxnet`.

To install MXNet on your system, you can follow the [MXNet documentation](https://mxnet.incubator.apache.org/versions/master/install/), or refer to our [scripts](scripts) or the `LIBRARY INSTALLATION` section in the [dockefiles](dockerfiles). OpenBLAS is used in our default build.

- The default blas is OpenBLAS.
  The default OpenBLAS path for macOS is `/usr/local/opt/openblas` if installed throught homebrew (openblas is keg-only, which means it was not symlinked into /usr/local, because macOS provides BLAS and LAPACK in the Accelerate framework).

- The default mxnet installation path is `/opt/mxnet` for linux, darwin and ppc64le without powerai; `/opt/DL/mxnet` for ppc64le with powerai.

- The default CUDA path is `/usr/local/cuda`

See [lib.go](lib.go) for details.

If you get an error about not being able to write to `/opt` then perform the following

```
sudo mkdir -p /opt/mxnet
sudo chown -R `whoami` /opt/mxnet
```

If you are using MXNet docker images or other libary paths, change CGO_CFLAGS, CGO_CXXFLAGS and CGO_LDFLAGS enviroment variables. Refer to [Using cgo with the go command](https://golang.org/cmd/cgo/#hdr-Using_cgo_with_the_go_command).

For example,

```
    export CGO_CFLAGS="${CGO_CFLAGS} -I/tmp/mxnet/include"
    export CGO_CXXFLAGS="${CGO_CXXFLAGS} -I/tmp/mxnet/include"
    export CGO_LDFLAGS="${CGO_LDFLAGS} -L/tmp/mxnet/lib"
```

### Go Packages

You can install the dependency through `go get`.

```
cd $GOPATH/src/github.com/c3sr/go-mxnet
go get -u -v ./...
```

Or use [Dep](https://github.com/golang/dep).

```
dep ensure -v
```

This installs the dependency in `vendor/`.

### Configure Environmental Variables

Configure the linker environmental variables since the MXNet C library is under a non-system directory. Place the following in either your `~/.bashrc` or `~/.zshrc` file

Linux
```
export LIBRARY_PATH=$LIBRARY_PATH:/opt/mxnet/lib
export LD_LIBRARY_PATH=/opt/mxnet/lib:$DYLD_LIBRARY_PATH

```

macOS
```
export LIBRARY_PATH=$LIBRARY_PATH:/opt/mxnet/lib
export DYLD_LIBRARY_PATH=/opt/mxnet/lib:$DYLD_LIBRARY_PATH
```

## Check the Build

Run `go build` in to check the dependences installation and library paths set-up.
On linux, the default is to use GPU, if you don't have a GPU, do `go build -tags nogpu` instead of `go build`.

**_Note_** : The CGO interface passes go pointers to the C API. This is an error by the CGO runtime. Disable the error by placing

```
export GODEBUG=cgocheck=0
```

in your `~/.bashrc` or `~/.zshrc` file and then run either `source ~/.bashrc` or `source ~/.zshrc`

## Examples

Examples of using the Go MXNet binding to do model inference are under [examples](examples).

### batch_mlmodelscope

This example shows how to use the MLModelScope tracer to profile the inference.

Refer to [Set up the external services](https://docs.mlmodelscope.org/installation/source/external_services/) to start the tracer.

Then run the example by

```
  cd example/batch_mlmodelscope
  go build
  ./batch
```

Now you can go to `localhost:16686` to look at the trace of that inference.

### batch_nvprof

This example shows how to use nvprof to profile the inference. You need GPU and CUDA to run this example.

```
  cd example/batch_nvprof
  go build
  nvprof --profile-from-start off ./batch_nvprof
```

Refer to [Profiler User's Guide](https://docs.nvidia.com/cuda/profiler-users-guide/index.html) for using nvprof.

module github.com/c3sr/go-mxnet

go 1.15

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
	github.com/jaegertracing/jaeger => github.com/uber/jaeger v1.22.0
	github.com/uber/jaeger => github.com/jaegertracing/jaeger v1.22.0
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)

require (
	github.com/anthonynsimon/bild v0.13.0
	github.com/c3sr/config v1.0.1
	github.com/c3sr/dlframework v1.0.2
	github.com/c3sr/downloadmanager v1.0.0
	github.com/c3sr/go-cupti v1.0.1
	github.com/c3sr/image v1.0.0
	github.com/c3sr/logger v1.0.1
	github.com/c3sr/nvidia-smi v1.0.0
	github.com/c3sr/tracer v1.0.0
	github.com/chewxy/math32 v1.0.6
	github.com/k0kubun/pp/v3 v3.0.7
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/unknwon/com v1.0.1
	gonum.org/v1/gonum v0.9.0
	gorgonia.org/tensor v0.9.14
)

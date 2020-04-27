module github.com/Basic-Components/jwttools

require (
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/etcd-io/etcd v3.3.20+incompatible
	github.com/golang/protobuf v1.3.3
	github.com/json-iterator/go v1.1.9
	github.com/liyue201/grpc-lb v0.0.0-20200224022010-2d7df6fc9b21
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/small-tk/pathlib v0.0.0-20190601032836-742166d9b695
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	google.golang.org/grpc v1.26.0
	gopkg.in/resty.v1 v1.12.0 // indirect
)

go 1.12

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.29.1

module github.com/dumacp/go-driverconsole

go 1.19

//github.com/AsynkronIT/protoactor-go v0.0.0-20220121183416-233df622d732
require (
	github.com/dumacp/go-logs v0.0.1
	github.com/dumacp/matrixorbital v0.0.0-20211112030057-e5299bc41d1f
	github.com/dumacp/pubsub v0.0.0-20200115200904-f16f29d84ee0
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/google/uuid v1.3.0 // indirect
	github.com/looplab/fsm v1.0.1
	go.etcd.io/bbolt v1.3.7
)

require (
	github.com/asynkron/protoactor-go v0.0.0-20230414121700-22ab527f4f7a
	github.com/dumacp/go-actors v0.0.0-20230503160549-734b3c336394
	github.com/dumacp/go-gwiot v0.0.0-00010101000000-000000000000
	github.com/dumacp/go-itinerary v0.0.0-20230427203726-7dd05dd6a3b5
	github.com/dumacp/go-levis v0.0.0-20230414205412-110e9cea515c
	github.com/dumacp/go-params v0.0.0-00010101000000-000000000000
	github.com/dumacp/go-schservices v0.0.1
	github.com/dumacp/gpsnmea v0.0.0-20201110195359-2994f05cfb52
	github.com/golang/geo v0.0.0-20210211234256-740aa86cb551
)

require (
	github.com/Workiva/go-datastructures v1.0.53 // indirect
	github.com/asynkron/gofun v0.0.0-20220329210725-34fed760f4c2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goburrow/modbus v0.1.0 // indirect
	github.com/goburrow/serial v0.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/lithammer/shortuuid/v4 v4.0.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/orcaman/concurrent-map v1.0.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	github.com/twmb/murmur3 v1.1.6 // indirect
	go.opentelemetry.io/otel v1.12.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.35.0 // indirect
	go.opentelemetry.io/otel/metric v0.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.12.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.12.0 // indirect
	golang.org/x/exp v0.0.0-20221012134508-3640c57a48ea // indirect
	golang.org/x/net v0.6.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/grpc v1.52.3 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/dumacp/go-fareCollection => ../go-fareCollection

replace github.com/dumacp/go-levis => ../go-levis

replace github.com/dumacp/go-params => ../go-params

replace github.com/dumacp/go-itinerary => ../go-itinerary

replace github.com/dumacp/go-gwiot => ../go-gwiot

replace github.com/dumacp/go-schservices => ../go-schservices

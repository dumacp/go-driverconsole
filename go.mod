module github.com/dumacp/go-driverconsole

go 1.21.0

toolchain go1.22.0

//github.com/AsynkronIT/protoactor-go v0.0.0-20220121183416-233df622d732
require (
	github.com/dumacp/go-logs v0.0.2-0.20241119230451-ac5d49be15ce
	github.com/dumacp/matrixorbital v0.0.0-20230818142648-d37611aecca5
	github.com/dumacp/pubsub v0.0.0-20200115200904-f16f29d84ee0
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/google/uuid v1.6.0
	github.com/looplab/fsm v1.0.2
	go.etcd.io/bbolt v1.3.10
)

require (
	github.com/asynkron/protoactor-go v0.0.0-20240413045429-76c172a71a16
	github.com/dumacp/go-actors v0.0.0-20240613144007-fcb38ee7b9b1
	github.com/dumacp/go-fareCollection v0.0.0-00010101000000-000000000000
	github.com/dumacp/go-gwiot v0.0.0-20250219205658-c5ffa6d680c1
	github.com/dumacp/go-ignition v0.0.0-20240301165217-62b8949edaf7
	github.com/dumacp/go-itinerary v0.0.0-20250206143001-1ce6638f2c19
	github.com/dumacp/go-levis v0.0.0-20241119224207-a91cacdf55e3
	github.com/dumacp/go-params v0.0.0-20250108191046-36f8cb3a96ac
	github.com/dumacp/go-schservices v0.0.4-0.20250115134655-33531cf6227c
	github.com/dumacp/gpsnmea v0.0.0-20201110195359-2994f05cfb52
	github.com/golang/geo v0.0.0-20210211234256-740aa86cb551
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/Workiva/go-datastructures v1.1.3 // indirect
	github.com/asynkron/gofun v0.0.0-20220329210725-34fed760f4c2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/dumacp/keycloak v0.0.0-20191212174805-9e9a5c3da24f // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goburrow/modbus v0.1.0 // indirect
	github.com/goburrow/serial v0.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lithammer/shortuuid/v4 v4.0.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/nats-io/nats.go v1.36.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/orcaman/concurrent-map v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/oauth2 v0.13.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/grpc v1.60.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

replace github.com/dumacp/go-fareCollection => ../go-fareCollection

//replace github.com/dumacp/go-levis => ../go-levis

//replace github.com/dumacp/go-params => ../go-params

//replace github.com/dumacp/go-itinerary => ../go-itinerary

replace github.com/dumacp/go-gwiot => ../go-gwiot

replace github.com/dumacp/go-schservices => ../go-schservices

//replace github.com/dumacp/go-actors => ../go-actors

//replace github.com/dumacp/matrixorbital => ../matrixorbital

replace github.com/asynkron/protoactor-go => ../../asynkron/protoactor-go

replace github.com/nats-io/nats.go => ../../nats-io/nats.go

//replace github.com/dumacp/go-logs => ../go-logs

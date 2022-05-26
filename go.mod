module github.com/dumacp/go-driverconsole

go 1.16

replace github.com/dumacp/go-levis => ../go-levis

replace github.com/dumacp/matrixorbital => ../matrixorbital

replace github.com/dumacp/go-logs => ../go-logs

//replace github.com/dumacp/go-levis/pkg/levis => ./pkg/levis

//github.com/AsynkronIT/protoactor-go v0.0.0-20220121183416-233df622d732
require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20220214042420-fcde2cd4013e
	//github.com/AsynkronIT/protoactor-go v0.0.0-20211124041449-becb6dbfc022
	github.com/dumacp/go-levis v0.0.0-20220524154038-12423a7dea34
	github.com/dumacp/go-logs v0.0.0-20220502162726-a75aa8a855e9
	github.com/dumacp/matrixorbital v0.0.0-20211112030057-e5299bc41d1f
	github.com/gogo/protobuf v1.3.2
	github.com/looplab/fsm v0.3.0
	github.com/orcaman/concurrent-map v1.0.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	golang.org/x/net v0.0.0-20220526153639-5463443f8c37 // indirect
	google.golang.org/genproto v0.0.0-20220526192754-51939a95c655 // indirect
)

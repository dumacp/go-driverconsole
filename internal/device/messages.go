package device

type MsgDevice struct {
	Device interface{}
}
type StartDevice struct{}
type StopDevice struct{}
type CloseDevice struct{}
type Subscribe struct{}

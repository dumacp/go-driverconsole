package device

type Device interface {
	Init() (interface{}, error)
	Close() error
}

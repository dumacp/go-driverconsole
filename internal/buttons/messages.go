package buttons

type MsgFatal struct{}
type MsgDevice struct {
	Device interface{}
}
type MsgInputText struct {
	Text string
}
type MsgSelectPaso struct{}
type MsgEnterPaso struct{}
type MsgMainScreen struct{}
type MsgEnterRuta struct {
	Route int
}
type MsgChangeRuta struct{}
type MsgConfirmation struct{}
type MsgWarning struct{}
type MsgResetCounter struct{}
type MsgInitRecorrido struct{}
type MsgStopRecorrido struct{}
type MsgSubscribe struct{}
type MsgMemory struct {
	Key   string
	Value interface{}
}

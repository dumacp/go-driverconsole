package buttons

type MsgFatal struct{}

type MsgInputText struct {
	Text string
}
type MsgSelectPaso struct{}
type MsgEnterPaso struct{}
type MsgDisableEnterPaso struct{}
type MsgEnableEnterPaso struct{}
type MsgMainScreen struct{}
type MsgEnterDriver struct {
	Driver int
}
type MsgEnterRuta struct {
	Route int
}
type MsgChangeRuta struct{}
type MsgChangeDriver struct{}
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
type MsgDeviceError struct{}
type MsgShowAlarms struct{}
type MsgReturnFromAlarms struct{}
type MsgReturnFromVehicle struct{}

type MsgBrightnessPlus struct{}
type MsgBrightnessMinus struct{}

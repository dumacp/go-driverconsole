package gps

type CounterMap struct {
}

type MsgSubscribe struct{}
type MsgRequestStatus struct {
}
type MsgGpsData struct {
	Data GPSData
}

type MsgGpsStatusRequest struct{}
type MsgGpsStatus struct {
	State bool
}

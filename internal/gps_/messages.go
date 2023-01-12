package gps

type CounterMap struct {
}

type MsgSubscribe struct{}
type MsgRequestStatus struct {
}
type MsgGpsData struct {
	Data GPSData
}

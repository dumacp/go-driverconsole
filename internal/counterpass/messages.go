package counterpass

type CounterMap struct {
	Inputs0    int64 `json:"inputs0"`
	Inputs1    int64 `json:"inputs1"`
	Outputs0   int64 `json:"outputs0"`
	Outputs1   int64 `json:"outputs1"`
	Anomalies0 int64 `json:"anomalies0"`
	Anomalies1 int64 `json:"anomalies1"`
	Tampering0 int64 `json:"tampering0"`
	Tampering1 int64 `json:"tampering1"`
}

type CounterEvent struct {
	Inputs  int
	Outputs int
}

type CounterExtraEvent struct {
	Text []byte
}

type MsgSubscribe struct{}
type MsgRequestStatus struct {
}
type MsgStatus struct {
	State bool
}
type TurnstileRegisters struct {
	Registers []uint32 `json:"registers"`
}

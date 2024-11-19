package app

type ValidationData struct {
	CountInputs  int32 `json:"countInputs"`
	CountOutputs int32 `json:"countOutputs"`
	CashInputs   int32 `json:"cashInputs"`
	ElectInputs  int32 `json:"electInputs"`
	Time         int64 `json:"timestamp"`
}

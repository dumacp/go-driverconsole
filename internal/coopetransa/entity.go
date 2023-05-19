package app

type ValidationData struct {
	CountInputs  int32 `json:"countInputs"`
	CountOutputs int32 `json:"countOutputs"`
	Time         int64 `json:"timestamp"`
}

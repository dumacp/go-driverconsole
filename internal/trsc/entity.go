package app

import (
	"fmt"
	"time"

	"github.com/dumacp/go-schservices/api/services"
)

type ValidationData struct {
	CountInputs  int32 `json:"countInputs"`
	CountOutputs int32 `json:"countOutputs"`
	CashInputs   int32 `json:"cashInputs"`
	ElectInputs  int32 `json:"electInputs"`
	Time         int64 `json:"timestamp"`
}

func (v *ValidationData) String() string {
	return fmt.Sprintf(`{"countInputs": %d, "countOutputs": %d, "cashInputs": %d, "electInputs": %d, "timestamp": %q}`,
		v.CountInputs, v.CountOutputs, v.CashInputs, v.ElectInputs, time.Unix(v.Time/1000, 0))
}

type CompanySchService struct {
	String       string
	ResumeString string
	Services     *services.ScheduleService
}

package platform

// DriverDailyServicesResponse is the top-level response from
// GET /api/external-system-gateway/rest/driver-daily-services/{document}
type DriverDailyServicesResponse struct {
	DriverID   string                    `json:"driverId"`
	Date       string                    `json:"date"`
	Services   []DriverDailyService      `json:"services"`
	Totals     DriverDailyServicesTotals `json:"totals"`
	Pagination DriverDailyPagination     `json:"pagination"`
}

// DriverDailyService represents one service row executed by the driver in the day.
type DriverDailyService struct {
	ID           string                 `json:"id"`
	VehiclePlate string                 `json:"vehiclePlate"`
	StartTime    int64                  `json:"startTime"` // epoch milliseconds
	EndTime      int64                  `json:"endTime"`   // epoch milliseconds
	RouteID      string                 `json:"routeId"`
	RouteName    string                 `json:"routeName"`
	FrontDoor    DriverDailyDoorCounts  `json:"frontDoor"`
	BackDoor     DriverDailyDoorCounts  `json:"backDoor"`
}

// DriverDailyDoorCounts holds inputs and outputs for a single door.
type DriverDailyDoorCounts struct {
	Inputs  int `json:"inputs"`
	Outputs int `json:"outputs"`
}

// DriverDailyServicesTotals holds aggregated counts.
type DriverDailyServicesTotals struct {
	TotalFrontDoorInputs  int `json:"totalFrontDoorInputs"`
	TotalFrontDoorOutputs int `json:"totalFrontDoorOutputs"`
	TotalBackDoorInputs   int `json:"totalBackDoorInputs"`
	TotalBackDoorOutputs  int `json:"totalBackDoorOutputs"`
}

// DriverDailyPagination holds pagination metadata.
type DriverDailyPagination struct {
	TotalServices int `json:"totalServices"`
	TotalPages    int `json:"totalPages"`
}

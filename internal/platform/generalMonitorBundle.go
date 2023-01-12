package platform

type UserObjListHumanResources struct {
	IID        string `json:"_id"`
	Active     bool   `json:"active"`
	DocumentID string `json:"documentId"`
	FirstName  string `json:"firstName"`
	ID         string `json:"id"`
	LastName   string `json:"lastName"`
}

type DriverBundle struct {
	ID         string `json:"id"`
	DocumentID string `json:"documentId"`
	FullName   string `json:"fullName"`
}

type RouteBundle struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	AuthorityCode string `json:"authorityCode"`
	Name          string `json:"name"`
	DivipolCode   string `json:"divipolCode"`
}

type HumanResources struct {
	DriversIds  []string                     `json:"driversIds"`
	ManagerIds  []string                     `json:"managerIds"`
	UserObjList []*UserObjListHumanResources `json:"userObjList"`
}

type VehicleBundle struct {
	Active   bool   `json:"active"`
	DeviceID string `json:"deviceId"`
	// HumanResources *HumanResources `json:"humanResources"`
	Id             string `json:"id"`
	InternalNumber string `json:"internalNumber"`
	Model          string `json:"model"`
	Plate          string `json:"plate"`
	Ttype          string `json:"type"`
	Year           int    `json:"year"`
}

type ItineraryBundle struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	FlatPath *FlatPath `json:"flatPath"`
}

type FlatPath struct {
	CheckPoints []*CheckPoint `json:"checkPoints"`
}

type CheckPoint struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Radius   int    `json:"radius"`
	Distance int    `json:"distance"`
	Eta      int    `json:"eta"`
}

type ServiceBundle struct {
	ID               string           `json:"id"`
	ScheduleDateTime int64            `json:"scheduleDateTime"`
	Route            *Route           `json:"route"`
	Driver           *DriverBundle    `json:"driver"`
	Vehicle          *VehicleBundle   `json:"vehicle"`
	Itinerary        *ItineraryBundle `json:"itinerary"`
}

type GeneralMonitorBundle struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Active   bool             `json:"active"`
	Vehicles []*VehicleBundle `json:"vehicles"`
	Services []*ServiceBundle `json:"services"`
}

type DataBundle struct {
	GeneralMonitorBundle *GeneralMonitorBundle `json:"GeneralMonitorBundle"`
}

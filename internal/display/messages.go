package display

type Device struct {
	Device interface{}
}
type Reset struct {
	// CountManual  int
	// CountParcial int
	// CountAppFare int
	// Route        string
	// Itininerary  string
}
type UpdateMainScreen struct {
	// timeLapse    int
	// CountManual  int
	// CountParcial int
	// CountAppFare int
	// Route        string
	// Itininerary  string
}
type UpdateDate struct{}
type DisplayCount struct {
	CountManual  int
	CountParcial int
	CountAppFare int
}
type DisplayError struct {
	Data string
}
type Route struct {
	Route       string
	Itininerary string
}
type Itininerary struct {
	Data string
}
type TimeLapse struct {
	Data int
}
type StopRecorrido struct{}
type InitRecorrido struct{}
type ResetCounter struct{}
type UpVoc struct{}
type EnterVoc struct{}
type DisableSelectPASO struct{}
type EnterPASO struct{}
type AddText struct {
	Text string
}
type DelText struct {
	Count int
}
type EnterText struct{}

type MsgRoutes struct {
	Data map[int]string
}

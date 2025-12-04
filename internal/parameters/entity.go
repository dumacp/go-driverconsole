package parameters

type TerminalConfig struct {
	Url                     string `json:"url"`
	TerminalPort            string `json:"port"`
	TerminalBaud            int    `json:"baud"`
	IsCashEnabled           bool   `json:"cash"`
	IsItineraryProgEnabled  bool   `json:"itinerary_programmed"`
	IsLegacy                bool   `json:"legacy"`
	IsReverseTQ             bool   `json:"reverse_tq"`
	IsEnableCameraFrontDoor bool   `json:"camera_frontdoor"`
}

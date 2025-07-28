package parameters

type TerminalConfig struct {
	TerminalPort           string `json:"port"`
	TerminalBaud           int    `json:"baud"`
	IsCashEnabled          bool   `json:"cash"`
	IsItineraryProgEnabled bool   `json:"itinerary_programmed"`
	IsLegacy               bool   `json:"legacy"`
}

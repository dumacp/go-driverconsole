package ui

type RegisterAddress int

const (
	AddrSelectPaso         RegisterAddress = 0
	AddrEnterPaso                          = 1
	AddrEnterRuta                          = 2
	AddrEnterDriver                        = 3
	AddrScreenAlarms                       = 7
	AddrScreenProgVeh                      = 8
	AddrScreenProgDriver                   = 9
	AddrScreenMore                         = 10
	AddrLedGsp                             = 12
	AddrLedNetwork                         = 11
	AddrAddBright                          = 21
	AddrSubBright                          = 22
	AddrTextDate                           = 60
	AddrTextRoute                          = 100
	AddrTextDriver                         = 140
	AddrNumInputs                          = 80
	AddrNumOutputs                         = 84
	AddrNumDeviation                       = 88
	AddrTextCurrentService                 = 160
	AddrLedCurrentService                  = 200
	AddrTextPopup                          = 400
	AddrTextConfirmation                   = 600
	AddrToggleConfirmation                 = 5
	AddrTogglePopup                        = 6
	AddrTextProgVeh                        = 2000
)

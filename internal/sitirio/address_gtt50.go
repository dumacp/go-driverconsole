package app

const (
	AddrGttSelectPaso      int = 5
	AddrGttScreenMain          = 1
	AddrGttScreenVeh           = 9
	AddrGttScreenDriver        = 6
	AddrGttScreenStats         = 8
	AddrGttScreenAlarm         = 5
	AddrGttScreenKeyRoute      = 2
	AddrGttScreenKeyDriver     = 3
	AddrGttScreenWarn          = 5
	AddrGttScreenInfo          = 4
	AddrGttEnterPaso           = 20

	AddrGttScreenSwitch          = -1
	AddrGttTogglePopup           = 17
	AddrGttToggleWarnPopup       = 33
	AddrGttEnterScreenAlarms     = 24
	AddrGttEnterScreenProgVeh    = 35
	AddrGttEnterScreenProgDriver = 28
	AddrGttEnterScreenMore       = 25
	AddrGttLedGsp                = 34
	AddrGttLedNetwork            = 88
	AddrGttAddBright             = 82
	AddrGttSubBright             = 89
	AddrGttLedBeep               = 23
	AddrGttSwitchStep            = 5
	AddrGttSendStep              = 20
	AddrGttTextDate              = 90
	AddrGttNumRoute              = 52
	AddrGttTextRoute             = 18
	AddrGttTextDriver            = 21
	AddrGttNumDriver             = 48
	AddrGttNumCashInputs         = 16
	AddrGttNumElectonicInputs    = 15
	AddrGttCounterInputs         = 3
	AddrGttCounterOutputs        = 14
	AddrGttNumDeviation          = -1
	AddrGttTextCurrentService    = -1
	AddrGttLedCurrentService     = -1
	AddrGttTextWarnPopup         = 33
	AddrGttTextPopup             = 17
	AddrGttTextWarning           = 29
	AddrGttReturnWarning         = 85
	AddrGttTextConfirmation      = 70
	AddrGttReturnConfimration    = 75
	AddrGttTextProgVeh           = 76
	AddrGttReturnProgVeh         = 72
	AddrGttTextProgDriver        = 78
	AddrGttReturnProgDriver      = 108
	AddrGttTextNotiAlarm         = 87
	AddrGttReturnNotiAlarm       = 84
	AddrGttLedReset              = -1
	AddrGttReturnStats           = 144
)

// Buttons Keyboard
const (
	AddrGttButtonRoute        = 77
	AddrGttButtonRoute_SEND   = 43
	AddrGttButtonRoute_Clear  = 42
	AddrGttButtonRoute_Cancel = 44
	AddrGttNoRoute            = 13
	AddrGttNameRoute          = 31
	AddrGttButtonRoute_0      = 40
	AddrGttButtonRoute_1      = 9
	AddrGttButtonRoute_2      = 19
	AddrGttButtonRoute_3      = 22
	AddrGttButtonRoute_4      = 27
	AddrGttButtonRoute_5      = 30
	AddrGttButtonRoute_6      = 31
	AddrGttButtonRoute_7      = 36
	AddrGttButtonRoute_8      = 37
	AddrGttButtonRoute_9      = 38

	AddrGttButtonDriver        = 81
	AddrGttButtonDriver_SEND   = 66
	AddrGttButtonDriver_Clear  = 65
	AddrGttButtonDriver_Cancel = 67
	AddrGttNoDriver            = 13
	AddrGttNameDriver          = 31
	AddrGttButtonDriver_0      = 61
	AddrGttButtonDriver_1      = 51
	AddrGttButtonDriver_2      = 53
	AddrGttButtonDriver_3      = 54
	AddrGttButtonDriver_4      = 55
	AddrGttButtonDriver_5      = 56
	AddrGttButtonDriver_6      = 57
	AddrGttButtonDriver_7      = 58
	AddrGttButtonDriver_8      = 59
	AddrGttButtonDriver_9      = 60
)

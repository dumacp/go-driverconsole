package buttons

type Button struct {
	Addr  int
	Value int
}

type ButtonDevice interface {
	ListenButtons() chan *Button
}

type ButtonValue int

const (
	TextNumLabel = "textNumLabel"
)

const (
	GridButton0 ButtonValue = iota
	GridButton1
	GridButton2
	GridButton3
	GridButton4
	GridButton5
	GridButton6
	GridButton7
	GridButton8
	GridButton9
	GridButtonEnter
	GridButtonDel
	UpVoc
	EnterVoc
	SelectPASO
	EnterPASO
	ResetCounter
	ResetRecorrido
)

// func ListenButtons(v *gtt43a.Event, getStateSelectPaso func() (byte, error)) ButtonValue {
// 	if v.Type == gtt43a.ButtonClick && v.Value[0] == 0x00 {
// 		switch int(v.ObjId) {
// 		case buttonEnter:
// 			return EnterVoc
// 		case buttonUp:
// 			return UpVoc
// 		case buttonEnterPaso:
// 			return EnterPASO
// 		case buttonRecorrido:
// 			return ResetRecorrido
// 		case buttonCounter:
// 			return ResetCounter
// 		case buttonSelectPaso:
// 			if state, err := getStateSelectPaso(); err == nil {
// 				log.Printf("state %X: [%X]\n", v.ObjId, state)
// 				if state == 0x01 {
// 					log.Println("SelectPASO")
// 					return SelectPASO
// 				}
// 			} else {
// 				log.Println(err)
// 			}
// 		}
// 	} else if v.Type == gtt43a.RegionTouch && v.Value[0] == 0x01 {
// 		switch int(v.ObjId) {
// 		case buttonGrid1:
// 			return GridButton1
// 		case buttonGrid2:
// 			return GridButton2
// 		case buttonGrid3:
// 			return GridButton3
// 		case buttonGrid4:
// 			return GridButton4
// 		case buttonGrid5:
// 			return GridButton5
// 		case buttonGrid6:
// 			return GridButton6
// 		case buttonGrid7:
// 			return GridButton7
// 		case buttonGrid8:
// 			return GridButton8
// 		case buttonGrid9:
// 			return GridButton9
// 		case buttonGrid0:
// 			return GridButton0
// 		case buttonGridEnter:
// 			return GridButtonEnter
// 		case buttonGridDel:
// 			return GridButtonDel
// 		}
// 	}
// 	return 0
// }

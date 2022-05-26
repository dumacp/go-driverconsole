//+build gtt43 !levis

package buttons

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/matrixorbital/gtt43a"
)

// type Button struct {
// 	Addr  int
// 	Value int
// }

// type ButtonDevice interface {
// 	ListenButtons() chan *Button
// }

// const (
// 	textCivica       int = 25
// 	textEfectivo     int = 27
// 	labelParcial     int = 5
// 	textParcial      int = 9
// 	usosSliceValue   int = 3
// 	labelError       int = 20
// 	labelTextInput   int = 29
// 	tittleRuta       int = 2
// 	textRuta         int = 1
// 	buttonEnter      int = 18
// 	buttonUp         int = 19
// 	buttonSelectPaso int = 15
// 	buttonEnterPaso  int = 16
// 	buttonRecorrido  int = 10
// 	buttonCounter    int = 6
// 	buttonGrid1      int = 0
// 	buttonGrid2      int = 3
// 	buttonGrid3      int = 6
// 	buttonGrid4      int = 1
// 	buttonGrid5      int = 4
// 	buttonGrid6      int = 7
// 	buttonGrid7      int = 2
// 	buttonGrid8      int = 5
// 	buttonGrid9      int = 8
// 	buttonGrid0      int = 10
// 	buttonGridEnter  int = 11
// 	buttonGridDel    int = 9
// 	timeRecorrido    int = 17
// 	timeHour         int = 24
// 	timeDate         int = 28
// 	textBoxError     int = 20
// )

// type ButtonValue int

const (
	addrSelectPaso     = 6
	addrEnterPaso      = 7
	addrEnterRuta      = 25
	addrChangeRuta     = 1
	addrResetRecorrido = 9
	addrConfirmation   = 33
	addrWarning        = 35

	addrDelRoute    = 24
	addrClearRoute  = 36
	addrCancelRoute = 22

	addrNoRoute   = 13
	addrNameRoute = 31
	addrButton_1  = 10
	addrButton_2  = 11
	addrButton_3  = 17
	addrButton_4  = 15
	addrButton_5  = 12
	addrButton_6  = 18
	addrButton_7  = 19
	addrButton_8  = 14
	addrButton_9  = 20
	addrButton_0  = 16
)

func DisableStep(dev interface{}) error {
	devv, ok := dev.(gtt43a.Display)
	if !ok {
		return fmt.Errorf("dev is not LEVIS device")
	}
	if err := devv.SetPropertyValueU8(addrEnterPaso, gtt43a.HasFocus)(0); err != nil {
		return err
	}
	return nil
}

func ListenButtons(dev interface{}, ctx actor.Context, mem <-chan *MsgMemory, quit <-chan int) error {

	devv, ok := dev.(gtt43a.Display)
	if !ok {
		return fmt.Errorf("dev is not LEVIS device")
	}

	if err := devv.Listen(); err != nil {
		logs.LogWarn.Println(err)
	}

	ch, err := devv.Events()
	if err != nil {
		return err
	}

	// rutaName := "RUTA"
	rutaCode := ""

	go func(ctx *actor.RootContext, self *actor.PID) {

		//defer close(ch)

		enableStep := time.NewTimer(5 * time.Second)
		activeStep := false
		valueResetRecorrido := 0

		for {
			select {
			case <-quit:
				return
			case <-enableStep.C:
				if activeStep {
					activeStep = false
					if err := devv.SetPropertyValueU8(addrEnterPaso, gtt43a.HasFocus)(0); err != nil {
						logs.LogWarn.Println(err)
					}
				}
			case v := <-mem:

				log.Printf("memory: %+v", v)
				switch v.Key {
				case TextNumLabel:
					rutaCode = fmt.Sprintf("%s", v.Value)
				}

			case v := <-ch:
				log.Printf("value (%d): [% X]", v.ObjId, v.Value)
				if len(v.Value) <= 0 {
					logs.LogWarn.Printf("value in event (%d) is invalid: %v", v.ObjId, v.Value)
					break
				}
				value := int(v.Value[0])
				switch int(v.ObjId) {
				case addrResetRecorrido:
					state, err := devv.GetPropertyValueU8(addrResetRecorrido, gtt43a.ButtonState)()
					if err != nil {
						logs.LogWarn.Println(err)
						break
					}
					log.Printf("state reset (%d): [%X]", addrResetRecorrido, state)
					if valueResetRecorrido != int(state) {
						if int(state) != 0 {
							ctx.Send(self, &MsgInitRecorrido{})
						} else {
							ctx.Send(self, &MsgStopRecorrido{})
						}
					}
				case addrChangeRuta:
					if value != 0 {

						ctx.Send(self, &MsgChangeRuta{})
					}
				case addrCancelRoute:
					if value != 0 {
						ctx.Send(self, &MsgMainScreen{})
					}
				case addrEnterRuta:
					if value != 0 {
						rutaCodeInt, _ := strconv.Atoi(rutaCode)
						ctx.Send(self, &MsgEnterRuta{Route: rutaCodeInt})
						rutaCode = ""
						// ctx.Send(self, &MsgMainScreen{})
					}
				case addrConfirmation:
					if value != 0 {
						data, err := devv.GetPropertyText(addrNoRoute, gtt43a.LabelText)()
						if err != nil {
							logs.LogWarn.Println(err)
							rutaCode = ""
						}
						log.Printf("keyNumLabel: %s", data)
						rutaCode = fmt.Sprintf("%s", data)
						ctx.Send(self, &MsgConfirmation{})
					}
				case addrWarning:
					if value != 0 {
						ctx.Send(self, &MsgWarning{})
					}
				case addrEnterPaso:
					if value == 0 {
						break
					}
					if !activeStep {
						break
					}
					activeStep = false
					if err := devv.SetPropertyValueU8(addrEnterPaso, gtt43a.HasFocus)(0); err != nil {
						logs.LogWarn.Println(err)
						break
					}
					ctx.Send(self, &MsgEnterPaso{})
				case addrSelectPaso:
					if value != 0 {
						// state, err := devv.GetPropertyValueU8(addrSelectPaso, gtt43a.ButtonState)()
						// if err != nil {
						// 	logs.LogWarn.Println(err)
						// 	break
						// }
						// log.Printf("state (%d): [%X]", addrSelectPaso, state)
						go func() {
							time.Sleep(300 * time.Millisecond)
							if err := devv.SetPropertyValueU8(addrSelectPaso, gtt43a.ButtonState)(0); err != nil {
								logs.LogWarn.Println(err)
							}
						}()
						err := devv.SetPropertyValueU8(addrEnterPaso, gtt43a.HasFocus)(1)
						if err != nil {
							logs.LogWarn.Println(err)
							break
						}
						if enableStep.Stop() {
							select {
							case <-enableStep.C:
							default:
							}
						}
						enableStep.Reset(5 * time.Second)
						activeStep = true
						ctx.Send(self, &MsgSelectPaso{})
					}
				case addrDelRoute:
					if value != 0 && len(rutaCode) > 0 {
						rutaCode = rutaCode[:len(rutaCode)-1]
						if err := devv.SetPropertyText(addrNoRoute, gtt43a.LabelText)(rutaCode); err != nil {
							logs.LogWarn.Println(err)
						}
					}
				case addrClearRoute:
					if value != 0 && len(rutaCode) > 0 {
						rutaCode = ""
					}
				case addrButton_0:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 0)
					}
				case addrButton_1:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 1)
					}
				case addrButton_2:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 2)
					}
				case addrButton_3:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 3)
					}
				case addrButton_4:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 4)
					}
				case addrButton_5:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 5)
					}
				case addrButton_6:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 6)
					}
				case addrButton_7:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 7)
					}
				case addrButton_8:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 8)
					}
				case addrButton_9:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 9)
					}
				}
			}
		}
	}(ctx.ActorSystem().Root, ctx.Self())
	return nil
}

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

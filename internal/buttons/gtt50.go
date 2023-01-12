//go:build (gtt50 || !levis) && (gtt50 || !gtt43)
// +build gtt50 !levis
// +build gtt50 !gtt43

package buttons

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/matrixorbital/gtt43a"
)

const (
	addrSelectPaso   = 5
	addrEnterPaso    = 20
	addrEnterRuta    = 43
	addrChangeRuta   = 77
	addrEnterDriver  = 66
	addrChangeDriver = 81
	// addrResetRecorrido = 9
	addrConfirmation = 75
	addrWarning      = 85

	addrClearRoute  = 42
	addrCancelRoute = 44

	addrNoRoute       = 13
	addrNameRoute     = 31
	addrButtonRoute_1 = 9
	addrButtonRoute_2 = 19
	addrButtonRoute_3 = 22
	addrButtonRoute_4 = 27
	addrButtonRoute_5 = 30
	addrButtonRoute_6 = 31
	addrButtonRoute_7 = 36
	addrButtonRoute_8 = 37
	addrButtonRoute_9 = 38
	addrButtonRoute_0 = 40

	addrClearDriver  = 65
	addrCancelDriver = 67
	addrButton_1     = 51
	addrButton_2     = 53
	addrButton_3     = 54
	addrButton_4     = 55
	addrButton_5     = 56
	addrButton_6     = 57
	addrButton_7     = 58
	addrButton_8     = 59
	addrButton_9     = 60
	addrButton_0     = 61

	addrShowAlarms        = 24
	addrReturnFromAlarms  = 84
	addrReturnFromVehicle = 144

	addrProgVehicle = 35
	addrProgDriver  = 28

	addrBrightnessPlus  = 82
	addrBrightnessMinus = 89
)

func EnableStep(dev interface{}) error {
	devv, ok := dev.(gtt43a.Display)
	if !ok {
		return fmt.Errorf("dev is not GTT device")
	}
	if err := devv.SetPropertyValueU8(addrEnterPaso, gtt43a.ButtonState)(0); err != nil {
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
	driverCode := ""

	go func(ctx *actor.RootContext, self *actor.PID) {

		//defer close(ch)

		lastStep := time.Now()
		enableStep := time.NewTimer(5 * time.Second)
		activeStep := false
		// valueResetRecorrido := 0

		for {
			select {
			case <-quit:
				return
			case <-enableStep.C:
				// fmt.Println("/////////////////  enableStep  ////////////////////")
				diff := time.Since(lastStep)
				if diff < 3*time.Second && activeStep {
					enableStep.Reset(diff)
					break
				}
				if activeStep {
					fmt.Println("reset addrSelectPaso")
					activeStep = false
					// if err := devv.SetPropertyValueU8(addrSelectPaso, gtt43a.ButtonState)(2); err != nil {
					// 	logs.LogWarn.Println(err)
					// }
					if err := devv.SetPropertyValueU8(addrSelectPaso, gtt43a.ButtonState)(0); err != nil {
						logs.LogWarn.Println(err)
					}
					// if err := devv.SetPropertyValueU8(addrSelectPaso, gtt43a.HasFocus)(0); err != nil {
					// 	logs.LogWarn.Println(err)
					// }
				}
			case v := <-mem:

				log.Printf("memory: %+v", v)
				switch v.Key {
				case TextNumRoute:
					rutaCode = fmt.Sprintf("%s", v.Value)
				case TextNumDriver:
					driverCode = fmt.Sprintf("%s", v.Value)
				}

			case v, ok := <-ch:
				if !ok {
					ctx.Send(self, &MsgDeviceError{})
					return
				}
				log.Printf("value (%d): [% X]", v.ObjId, v.Value)
				if len(v.Value) <= 0 {
					log.Printf("value in event (%d) is invalid: %v", v.ObjId, v.Value)
					logs.LogWarn.Printf("value in event (%d) is invalid: %v", v.ObjId, v.Value)
					break
				}
				value := int(v.Value[0])
				switch int(v.ObjId) {
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
						fmt.Printf("routeCode: %s\n", rutaCode)
						rutaCodeInt, _ := strconv.Atoi(rutaCode)
						if rutaCodeInt == 0 && len(driverCode) <= 0 {
							ctx.Send(self, &MsgEnterRuta{Route: -1})
						} else {
							ctx.Send(self, &MsgEnterRuta{Route: rutaCodeInt})
						}
						rutaCode = ""
					}
				case addrChangeDriver:
					if value != 0 {
						ctx.Send(self, &MsgChangeDriver{})
					}
				case addrCancelDriver:
					if value != 0 {
						ctx.Send(self, &MsgMainScreen{})
					}
				case addrEnterDriver:
					if value != 0 {
						driverCodeInt, _ := strconv.Atoi(driverCode)
						if driverCodeInt == 0 && len(driverCode) <= 0 {
							ctx.Send(self, &MsgEnterDriver{Driver: -1})
						} else {
							ctx.Send(self, &MsgEnterDriver{Driver: driverCodeInt})
						}
						driverCode = ""
					}
				case addrConfirmation:
					if value != 0 {
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
					if time.Since(lastStep) < 300*time.Millisecond {
						break
					}
					lastStep = time.Now()
					ctx.Send(self, &MsgEnterPaso{})
					if err := devv.SetPropertyValueU8(addrEnterPaso, gtt43a.ButtonState)(2); err != nil {
						logs.LogWarn.Println(err)
					}
					go func() {
						<-time.After(2 * time.Second)
						if err := devv.SetPropertyValueU8(addrEnterPaso, gtt43a.ButtonState)(0); err != nil {
							logs.LogWarn.Println(err)
						}
					}()
				case addrSelectPaso:
					if value != 0 {

						// if !enableStep.Stop() {
						// 	select {
						// 	case <-enableStep.C:
						// 	default:
						// 	}
						// }

						if activeStep {
							activeStep = false

							// if err := devv.SetPropertyValueU8(addrSelectPaso, gtt43a.ButtonState)(2); err != nil {
							// 	logs.LogWarn.Println(err)
							// }
							if err := devv.SetPropertyValueU8(addrSelectPaso, gtt43a.ButtonState)(0); err != nil {
								logs.LogWarn.Println(err)
							}
							break
						}
						// if err := devv.SetPropertyValueU8(addrSelectPaso, gtt43a.HasFocus)(1); err != nil {
						// 	logs.LogWarn.Println(err)
						// }
						if err := devv.SetPropertyValueU8(addrSelectPaso, gtt43a.ButtonState)(1); err != nil {
							logs.LogWarn.Println(err)
						}
						if !enableStep.Stop() {
							select {
							case <-enableStep.C:
							default:
							}
						}
						enableStep.Reset(10 * time.Second)
						activeStep = true
						ctx.Send(self, &MsgSelectPaso{})
					}
				case addrShowAlarms:
					if value != 0 {
						ctx.Send(self, &MsgShowAlarms{})
					}
				case addrReturnFromAlarms:
					if value != 0 {
						ctx.Send(self, &MsgReturnFromAlarms{})
					}
				case addrReturnFromVehicle:
					if value != 0 {
						ctx.Send(self, &MsgReturnFromVehicle{})
					}
				case addrClearRoute:
					if value != 0 && len(rutaCode) > 0 {
						rutaCode = ""
					}
				case addrClearDriver:
					if value != 0 && len(driverCode) > 0 {
						driverCode = ""
					}
				case addrButtonRoute_0:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 0)
					}
				case addrButtonRoute_1:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 1)
					}
				case addrButtonRoute_2:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 2)
					}
				case addrButtonRoute_3:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 3)
					}
				case addrButtonRoute_4:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 4)
					}
				case addrButtonRoute_5:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 5)
					}
				case addrButtonRoute_6:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 6)
					}
				case addrButtonRoute_7:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 7)
					}
				case addrButtonRoute_8:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 8)
					}
				case addrButtonRoute_9:
					if value != 0 {
						rutaCode = fmt.Sprintf("%s%d", rutaCode, 9)
					}
				case addrButton_0:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 0)
					}
				case addrButton_1:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 1)
					}
				case addrButton_2:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 2)
					}
				case addrButton_3:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 3)
					}
				case addrButton_4:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 4)
					}
				case addrButton_5:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 5)
					}
				case addrButton_6:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 6)
					}
				case addrButton_7:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 7)
					}
				case addrButton_8:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 8)
					}
				case addrButton_9:
					if value != 0 {
						driverCode = fmt.Sprintf("%s%d", driverCode, 9)
					}
				case addrBrightnessMinus:
					if value != 0 {
						ctx.Send(self, &MsgBrightnessMinus{})
					}
				case addrBrightnessPlus:
					if value != 0 {
						ctx.Send(self, &MsgBrightnessPlus{})
					}
				}
			}
		}
	}(ctx.ActorSystem().Root, ctx.Self())
	return nil
}

package app

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/ui"
	"github.com/dumacp/go-schservices/api/services"
)

func (a *App) listProg(msg *ListProgVeh) error {
	if len(a.companySchServices) <= 0 {
		return nil
	}
	dataSlice := make([]string, 0)
	// sliceSvc := make([]*CompanySchService, 0)

	ss := make([]*services.ScheduleService, 0)
	for _, v := range a.companySchServices {
		ss = append(ss, v)
	}
	cs := make([]*CompanySchService, 0)

	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].GetScheduleDateTime() < ss[j].GetScheduleDateTime()
	})

	slices.Reverse(ss)

	fmt.Printf("services oldest: %s\n", time.UnixMilli(ss[len(ss)-1].GetScheduleDateTime()))
	fmt.Printf("services newest: %s\n", time.UnixMilli(ss[0].GetScheduleDateTime()))

	untilSlice := make([]*services.ScheduleService, 0)
	for _, v := range ss {
		ts := time.UnixMilli(v.GetScheduleDateTime())
		if time.Until(ts) < 0 {
			break
		}
		untilSlice = append(untilSlice, v)
	}
	fmt.Printf("services until: %d\n", len(untilSlice))

	if len(untilSlice) <= 0 {
		return fmt.Errorf("no hay servicios disponibles")
	}

	fmt.Printf("services until oldest: %s\n", time.UnixMilli(untilSlice[len(untilSlice)-1].GetScheduleDateTime()))
	fmt.Printf("services until newest: %s\n", time.UnixMilli(untilSlice[0].GetScheduleDateTime()))

	// a.companySchServicesShow = make([]*CompanySchService, 0)

	slices.Reverse(untilSlice)

	for _, v := range untilSlice {
		// v := untilSlice[len(untilSlice)-1-i]
		if msg.Itinerary > 0 && v.GetItinerary().GetId() != int32(msg.Itinerary) {
			continue
		}
		ts := time.UnixMilli(v.GetScheduleDateTime())
		data := strings.ToLower(fmt.Sprintf(" %s: %s (%s)", ts.Format("01/02 15:04"),
			v.GetItinerary().GetName(), v.GetRoute().GetCode()))
		svc := strings.ToLower(fmt.Sprintf(`  id: %q
  estado: %q
  tiempo de inicio: %s
  itinenario: %q (%d)
  ruta: %q (%d)`,
			v.GetId(), v.GetState(), time.UnixMilli(v.GetScheduleDateTime()).Format("01/02 15:04:05"),
			v.GetItinerary().GetName(), v.GetItinerary().GetId(),
			v.GetRoute().GetCode(), v.GetRoute().GetId()))
		fmt.Printf("servicio: %v\n", svc)
		fmt.Printf("data: %s\n", data)
		cs = append(cs, &CompanySchService{
			String:       svc,
			ResumeString: data,
			Services:     v,
		})
		if len(cs) >= 9 {
			break
		}
	}

	a.companySchServicesShow = make([]*CompanySchService, 0)
	if len(cs) > 0 {
		dataSlice = append(dataSlice, cs[0].ResumeString)
	}
	if len(cs) > 0 {
		dataSlice = append(dataSlice, cs[0].ResumeString)
		a.companySchServicesShow = append(a.companySchServicesShow, cs[0])
	}
	for i := range cs {
		dataSlice = append(dataSlice, cs[i].ResumeString)
		a.companySchServicesShow = append(a.companySchServicesShow, cs[i])
	}
	if len(dataSlice) < 9 {
		for i := 0; i <= 9-len(dataSlice); i++ {
			size := Label2DisplayRegister(ui.PROGRAMATION_VEH_SCREEN).Size
			// un string de tamaño size de espacios
			spaces := strings.Repeat(" ", size)
			dataSlice = append(dataSlice, spaces)
		}
	}

	fmt.Printf("dataslice: %v\n", dataSlice)
	if len(dataSlice) > 0 {
		if err := a.uix.ShowProgVeh(dataSlice...); err != nil {
			return fmt.Errorf("event ShowProgVeh error: %s", err)
		}
		if err := a.uix.SetLed(AddrUpdateDropListProgVeh, false); err != nil {
			return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
		}
		if err := a.uix.SetLed(AddrUpdateDropListProgVeh, true); err != nil {
			return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
		}
	} else {
		return fmt.Errorf(`no hay servicios dsiponibles
para el itinerario: %d`, msg.Itinerary)
	}
	return nil
}

func (a *App) listProgShitfs(msg *ListShiftsVeh) error {
	if len(a.CompanyShiftsService) <= 0 {
		return nil
	}
	dataSlice := make([]string, 0)
	// sliceSvc := make([]*CompanySchService, 0)

	ss := make([]*services.ShiftService, 0)
	for _, v := range a.CompanyShiftsService {
		ss = append(ss, v)
	}
	cs := make([]*CompanyShift, 0)

	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].GetNextServiceTimeStamp() < ss[j].GetNextServiceTimeStamp()
	})

	slices.Reverse(ss)

	fmt.Printf("service shift oldest: %s\n", time.UnixMilli(ss[len(ss)-1].GetNextServiceTimeStamp()))
	fmt.Printf("service shiff newest: %s\n", time.UnixMilli(ss[0].GetNextServiceTimeStamp()))

	untilSlice := make([]*services.ShiftService, 0)
	for _, v := range ss {
		ts := time.UnixMilli(v.GetNextServiceTimeStamp())
		if time.Until(ts) < 10*time.Minute {
			break
		}
		untilSlice = append(untilSlice, v)
	}
	fmt.Printf("services until: %d\n", len(untilSlice))

	if len(untilSlice) <= 0 {
		return fmt.Errorf("no hay servicios disponibles")
	}

	fmt.Printf("services until oldest: %s\n", time.UnixMilli(untilSlice[len(untilSlice)-1].GetNextServiceTimeStamp()))
	fmt.Printf("services until newest: %s\n", time.UnixMilli(untilSlice[0].GetNextServiceTimeStamp()))

	// a.companySchServicesShow = make([]*CompanySchService, 0)

	slices.Reverse(untilSlice)

	for _, v := range untilSlice {
		// v := untilSlice[len(untilSlice)-1-i]
		if v.GetItinerary() == nil {
			continue
		}
		if len(msg.Shift) > 0 && v.GetShift() != msg.Shift {
			continue
		}
		ts := time.UnixMilli(v.GetNextServiceTimeStamp())
		data := strings.ToLower(fmt.Sprintf(" %s: (%q) %s (%s)", ts.Format("01/02 15:04"),
			v.GetShift(), v.GetItinerary().GetName(), v.GetRoute().GetCode()))
		svc := strings.ToLower(fmt.Sprintf(`  shift id: %q
  inicio prox. servicio: %s
  itinenario: %q (%d)
  ruta: %q (%d)
  versión prog: %q`,
			v.GetShift(),
			time.UnixMilli(v.GetNextServiceTimeStamp()).Format("01/02 15:04:05"),
			v.GetItinerary().GetName(), v.GetItinerary().GetId(),
			v.GetRoute().GetCode(), v.GetRoute().GetId(),
			v.GetVersion()))
		fmt.Printf("servicio: %v\n", svc)
		fmt.Printf("data: %s\n", data)
		cs = append(cs, &CompanyShift{
			String:       svc,
			ResumeString: data,
			Shift:        v,
		})
		if len(cs) >= 9 {
			break
		}
	}

	a.companyShiftsShow = make([]*CompanyShift, 0)
	if len(cs) > 0 {
		dataSlice = append(dataSlice, cs[0].ResumeString)
	}
	if len(cs) > 0 {
		dataSlice = append(dataSlice, cs[0].ResumeString)
		a.companyShiftsShow = append(a.companyShiftsShow, cs[0])
	}
	for i := range cs {
		dataSlice = append(dataSlice, cs[i].ResumeString)
		a.companyShiftsShow = append(a.companyShiftsShow, cs[i])
	}
	if len(dataSlice) < 9 {
		for i := 0; i <= 9-len(dataSlice); i++ {
			size := Label2DisplayRegister(ui.PROGRAMATION_VEH_SCREEN).Size
			// un string de tamaño size de espacios
			spaces := strings.Repeat(" ", size)
			dataSlice = append(dataSlice, spaces)
		}
	}

	fmt.Printf("dataslice: %v\n", dataSlice)
	if len(dataSlice) > 0 {
		if err := a.uix.ShowProgVeh(dataSlice...); err != nil {
			return fmt.Errorf("event ShowProgVeh error: %s", err)
		}
		if err := a.uix.SetLed(AddrUpdateDropListProgVeh, false); err != nil {
			return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
		}
		if err := a.uix.SetLed(AddrUpdateDropListProgVeh, true); err != nil {
			return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
		}
	} else {
		return fmt.Errorf(`no hay servicios dsiponibles
para el turno: %s`, msg.Shift)
	}
	return nil
}

func (a *App) PrepSliceServices(cs []*CompanySchService) error {
	if len(cs) <= 0 {
		return fmt.Errorf("no hay servicios disponibles")
	}
	dataSlice := make([]string, 0)
	// sliceSvc := make([]*CompanySchService, 0)

	a.companySchServicesShow = make([]*CompanySchService, 0)
	if len(cs) > 0 {
		dataSlice = append(dataSlice, cs[0].ResumeString)
	}
	if len(cs) > 0 {
		dataSlice = append(dataSlice, cs[0].ResumeString)
		a.companySchServicesShow = append(a.companySchServicesShow, cs[0])
	}
	for i := range cs {
		dataSlice = append(dataSlice, cs[i].ResumeString)
		a.companySchServicesShow = append(a.companySchServicesShow, cs[i])
	}
	if len(dataSlice) < 9 {
		for i := 0; i <= 9-len(dataSlice); i++ {
			size := Label2DisplayRegister(ui.PROGRAMATION_VEH_SCREEN).Size
			// un string de tamaño size de espacios
			spaces := strings.Repeat(" ", size)
			dataSlice = append(dataSlice, spaces)
		}
	}

	fmt.Printf("dataslice: %v\n", dataSlice)
	if len(dataSlice) > 0 {
		if err := a.uix.ShowProgVeh(dataSlice...); err != nil {
			return fmt.Errorf("event ShowProgVeh error: %s", err)
		}
		if err := a.uix.SetLed(AddrUpdateDropListProgVeh, false); err != nil {
			return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
		}
		if err := a.uix.SetLed(AddrUpdateDropListProgVeh, true); err != nil {
			return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
		}
	} else {
		return fmt.Errorf("no hay servicios disponibles")
	}
	return nil
}

func (a *App) GenereCompanySchServices(msg *ListProgVeh) ([]*CompanySchService, error) {
	if len(a.companySchServices) <= 0 {
		return nil, fmt.Errorf("no hay servicios disponibles")
	}
	// dataSlice := make([]string, 0)
	// sliceSvc := make([]*CompanySchService, 0)

	ss := make([]*services.ScheduleService, 0)
	for _, v := range a.companySchServices {
		ss = append(ss, v)
	}
	cs := make([]*CompanySchService, 0)

	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].GetScheduleDateTime() < ss[j].GetScheduleDateTime()
	})

	slices.Reverse(ss)

	fmt.Printf("services oldest: %s\n", time.UnixMilli(ss[len(ss)-1].GetScheduleDateTime()))
	fmt.Printf("services newest: %s\n", time.UnixMilli(ss[0].GetScheduleDateTime()))

	untilSlice := make([]*services.ScheduleService, 0)
	for _, v := range ss {
		ts := time.UnixMilli(v.GetScheduleDateTime())
		if time.Until(ts) < 0 {
			break
		}
		untilSlice = append(untilSlice, v)
	}
	fmt.Printf("services until: %d\n", len(untilSlice))

	if len(untilSlice) <= 0 {
		return nil, fmt.Errorf("no hay servicios disponibles")
	}

	fmt.Printf("services until oldest: %s\n", time.UnixMilli(untilSlice[len(untilSlice)-1].GetScheduleDateTime()))
	fmt.Printf("services until newest: %s\n", time.UnixMilli(untilSlice[0].GetScheduleDateTime()))

	// a.companySchServicesShow = make([]*CompanySchService, 0)

	slices.Reverse(untilSlice)

	for _, v := range untilSlice {
		// v := untilSlice[len(untilSlice)-1-i]
		if msg.Itinerary > 0 && v.GetItinerary().GetId() != int32(msg.Itinerary) {
			continue
		}
		ts := time.UnixMilli(v.GetScheduleDateTime())
		data := strings.ToLower(fmt.Sprintf(" %s: %s (%s)", ts.Format("01/02 15:04"),
			v.GetItinerary().GetName(), v.GetRoute().GetCode()))
		svc := strings.ToLower(fmt.Sprintf(`  id: %q
  estado: %q
  tiempo de inicio: %s
  itinenario: %q (%d)
  ruta: %q (%d)`,
			v.GetId(), v.GetState(), time.UnixMilli(v.GetScheduleDateTime()).Format("01/02 15:04:05"),
			v.GetItinerary().GetName(), v.GetItinerary().GetId(),
			v.GetRoute().GetCode(), v.GetRoute().GetId()))
		fmt.Printf("servicio: %v\n", svc)
		fmt.Printf("data: %s\n", data)
		cs = append(cs, &CompanySchService{
			String:       svc,
			ResumeString: data,
			Services:     v,
		})
		if len(cs) >= 9 {
			break
		}
	}

	if len(cs) <= 0 {
		return nil, fmt.Errorf("no hay servicios dsiponibles")
	}

	return cs, nil
}

func (a *App) requestProg(ctx actor.Context, msg *RequestProgVeh) error {
	// sí hya servicios disponibles desde la última consulta
	fmt.Printf("cs: %d, route: %d, iti: %d\n", len(a.companySchServices), a.route, msg.Itinerary)
	if time.Since(a.lastReqProgVeh) < 5*time.Minute && len(a.companySchServices) > 0 && (a.route == 0 || a.route == msg.Itinerary) {
		ss := a.companySchServicesShow
		if len(ss) > 0 {
			fmt.Printf("time: %s\n", time.UnixMilli(ss[len(ss)-1].Services.GetScheduleDateTime()))
		}
		// if len(ss) > 0 && time.Since(time.UnixMilli(ss[len(ss)-1].Services.GetScheduleDateTime())) < 1*time.Hour {
		ctx.Send(ctx.Self(), &ListProgVeh{
			Itinerary: msg.Itinerary,
		})
		return nil
		// }
	}
	if a.pidSvc == nil {
		return fmt.Errorf("service pid is nil")
	}
	res, err := ctx.RequestFuture(a.pidSvc, &services.GetCompanyProgSvcMsg{
		RouteId:   int32(msg.Itinerary),
		State:     services.State_SCHEDULED.String(),
		CompanyId: a.companyId,
		DeviceId:  a.platformId,
	}, 6*time.Second).Result()
	if err != nil {
		return fmt.Errorf("request service error: %s", err)
	}
	switch rs := res.(type) {
	case *services.CompanyProgSvcMsg:
		if len(rs.GetError()) > 0 {
			return fmt.Errorf(rs.GetError())
		} else {
			a.lastReqProgVeh = time.Now()
			cs := make(map[string]*services.ScheduleService, 0)
			if len(rs.GetScheduledServices()) > 0 {
				for _, v := range rs.GetScheduledServices() {
					cs[v.GetId()] = v
				}
				if len(cs) > 0 {
					fmt.Printf("services in company (%q): %d\n", a.companyId, len(cs))
					a.companySchServices = cs
				}
			}
			ctx.Send(ctx.Self(), &ListProgVeh{
				Itinerary: msg.Itinerary,
			})
		}
	default:
		return fmt.Errorf("response type error: %T", rs)
	}
	return nil
}

func (a *App) listDriverProg(msg *ListProgDriver) error {
	if len(a.shcservices) <= 0 {
		fmt.Printf("no hay servicios disponibles ListProgDriver\n")
		return nil
	}
	dataSlice := make([]string, 0)
	// sliceSvc := make([]*CompanySchService, 0)

	ss := make([]*services.ScheduleService, 0)
	for _, v := range a.shcservices {
		ss = append(ss, v)
	}
	cs := make([]*CompanySchService, 0)

	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].GetScheduleDateTime() < ss[j].GetScheduleDateTime()
	})

	slices.Reverse(ss)

	untilSlice := make([]*services.ScheduleService, 0)
	for _, v := range ss {
		ts := time.UnixMilli(v.GetScheduleDateTime())
		if time.Until(ts) < 0 {
			break
		}
		untilSlice = append(untilSlice, v)
	}

	if len(untilSlice) <= 0 {
		return fmt.Errorf("no hay servicios disponibles ListProgDriver")
	}

	slices.Reverse(untilSlice)

	// a.companySchServicesShow = make([]*CompanySchService, 0)

	for _, v := range untilSlice {

		if msg.Itinerary > 0 && v.GetItinerary().GetId() != int32(msg.Itinerary) {
			continue
		}
		if len(msg.DriverDocument) > 0 && v.GetDriver().GetDocument() != msg.DriverDocument {
			continue
		}
		ts := time.UnixMilli(v.GetScheduleDateTime())
		data := strings.ToLower(fmt.Sprintf(" %s: %s (%s) %s", ts.Format("01/02 15:04"),
			v.GetItinerary().GetName(), v.GetRoute().GetCode(), v.GetState()))
		svc := strings.ToLower(fmt.Sprintf(`  id: %q
  estado: %q
  tiempo de inicio: %s
  itinenario: %q (%d)
  ruta: %q (%d)`,
			v.GetId(), v.GetState(), time.UnixMilli(v.GetScheduleDateTime()).Format("01/02 15:04:05"),
			v.GetItinerary().GetName(), v.GetItinerary().GetId(),
			v.GetRoute().GetCode(), v.GetRoute().GetId()))
		fmt.Printf("servicio: %v\n", svc)
		fmt.Printf("data: %s\n", data)
		cs = append(cs, &CompanySchService{
			String:       svc,
			ResumeString: data,
			Services:     v,
		})
		if len(cs) >= 12 {
			break
		}
	}

	a.vehicleSchServicesShow = make([]*CompanySchService, 0)

	for i := range cs {
		dataSlice = append(dataSlice, cs[len(cs)-i-1].ResumeString)
		a.vehicleSchServicesShow = append(a.vehicleSchServicesShow, cs[len(cs)-i-1])
	}
	if len(dataSlice) < 12 {
		for i := 0; i <= 12-len(dataSlice); i++ {
			size := Label2DisplayRegister(ui.PROGRAMATION_DRIVER_TEXT).Size
			// un string de tamaño size de espacios
			spaces := strings.Repeat(" ", size)
			dataSlice = append(dataSlice, spaces)
		}
	}

	fmt.Printf("dataslice progDriver: %v\n", dataSlice)
	if len(dataSlice) > 0 {
		if err := a.uix.ShowProgDriver(dataSlice...); err != nil {
			return fmt.Errorf("event ShowProgVeh error: %s", err)
		}
		// if err := a.uix.SetLed(AddrUpdateDropListProgVeh, false); err != nil {
		// 	return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
		// }
		// if err := a.uix.SetLed(AddrUpdateDropListProgVeh, true); err != nil {
		// 	return fmt.Errorf("error setLed (AddrScreenProgVeh): %s", err)
		// }
	}
	return nil
}

func (a *App) requestProgShifts(ctx actor.Context, msg *RequestShitfsVeh) error {
	// sí hya servicios disponibles desde la última consulta
	fmt.Printf("cs: %d, route: %d, shift: %s\n", len(a.companySchServices), a.route, msg.Shift)
	if time.Since(a.lastReqShifts) < 5*time.Minute && len(a.CompanyShiftsService) > 0 && (a.shift == nil || a.shift.Shift != msg.Shift) {
		ss := a.companySchServicesShow
		if len(ss) > 0 {
			fmt.Printf("time: %s\n", time.UnixMilli(ss[len(ss)-1].Services.GetScheduleDateTime()))
		}
		// if len(ss) > 0 && time.Since(time.UnixMilli(ss[len(ss)-1].Services.GetScheduleDateTime())) < 1*time.Hour {
		ctx.Send(ctx.Self(), &ListShiftsVeh{
			Shift: msg.Shift,
		})
		return nil
		// }
	}
	if a.pidSvc == nil {
		return fmt.Errorf("service pid is nil")
	}
	res, err := ctx.RequestFuture(a.pidSvc, &services.GetCompanyProgShiftsMsg{
		ShiftId:   msg.Shift,
		CompanyId: a.companyId,
		DeviceId:  a.platformId,
	}, 6*time.Second).Result()
	if err != nil {
		return fmt.Errorf("request shifts service error: %s", err)
	}
	switch rs := res.(type) {
	case *services.CompanyProgShiftsMsg:
		if len(rs.GetError()) > 0 {
			return fmt.Errorf(rs.GetError())
		} else {
			a.lastReqProgVeh = time.Now()
			cs := make(map[string]*services.ShiftService, 0)
			if len(rs.GetShiftsServices()) > 0 {
				for _, v := range rs.GetShiftsServices() {
					cs[v.Shift] = v
				}
				if len(cs) > 0 {
					fmt.Printf("shifts services in company (%q): %d\n", a.companyId, len(cs))
					a.CompanyShiftsService = cs
				}
			}
			ctx.Send(ctx.Self(), &ListShiftsVeh{
				Shift: msg.Shift,
			})
		}
	default:
		return fmt.Errorf("response shifts type error: %T", rs)
	}
	return nil
}

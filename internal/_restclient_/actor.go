package restclient

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/app"
	"github.com/dumacp/go-driverconsole/internal/itinerary"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/golang/geo/s2"
	"golang.org/x/oauth2"
)

const (
	serviceURL               = "%s/api/external-system-gateway/rest/device-service"
	itineraryURL             = "%s/api/external-system-gateway/rest/device_itineraries"
	bundleURL                = "%s/api/emi-gateway/graphql/http"
	defaultUsername          = "dev.nebulae"
	filterHttpQuery          = "?deviceId=%s&scheduledServices=%v&liveExecutedServices=%v"
	filterItineraryHttpQuery = "?page=%d&count=%d&active=true"
	defaultPassword          = "uno.2.tres"
)

const TIMEOUT = 3 * time.Minute

type Actor struct {
	lastGetService time.Time
	id             string
	userHttp       string
	passHttp       string
	clientid       string
	clientsecret   string
	realm          string
	url            string
	bundleUrl      string
	itiUrl         string
	keyUrl         string
	tks            oauth2.TokenSource
	httpClient     *http.Client
}

func NewActor(id, user, pass, url, keyUrl, realm, clientid, clientsecret string) actor.Actor {
	a := &Actor{}
	a.id = id
	a.userHttp = user
	a.passHttp = pass
	a.url = url
	a.keyUrl = keyUrl
	a.clientid = clientid
	a.clientsecret = clientsecret
	a.realm = realm
	return a
}

func (a *Actor) Receive(ctx actor.Context) {
	fmt.Printf("Message arrived in %s: %s, %T, %s\n",
		ctx.Self().GetId(), ctx.Message(), ctx.Message(), ctx.Sender())
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.itiUrl = fmt.Sprintf(itineraryURL, a.url)
		a.bundleUrl = fmt.Sprintf(bundleURL, a.url)

		logs.LogInfo.Printf("started \"%s\", %v", ctx.Self().GetId(), ctx.Self())

	case *actor.Stopping:

	case *app.MsgGetItinieary:
		if err := func() error {
			if time.Since(a.lastGetService) < 30*time.Second {
				return fmt.Errorf("last GetSchedule was before 30 seconds")
			}
			if a.tks == nil || !(func() bool {
				t, err := a.tks.Token()
				if err != nil {
					return false
				}
				return time.Until(t.Expiry) > 0
			}()) {
				tks, c, err := Token(a.userHttp, a.passHttp, a.keyUrl, a.realm, a.clientid, a.clientsecret)
				if err != nil {
					return err
				}
				a.tks = tks
				a.httpClient = c
			}
			if a.httpClient != nil {
				a.httpClient.CloseIdleConnections()
			}

			tk, _ := a.tks.Token()

			log.Printf("token: %s", tk)

			itiPlatform, _, err := PlataformRequestItinerary(a.httpClient, a.bundleUrl, msg.ID, msg.OrganizationID)
			if err != nil {
				a.httpClient = nil
				return err
			}

			controlPoints := make([]itinerary.ControlPoint, 0)
			for _, v := range itiPlatform.Path {
				coord := v.Coords
				if len(coord) < 2 {
					continue
				}
				latlon := s2.LatLngFromDegrees(coord[1], coord[0])
				point := s2.PointFromLatLng(latlon)
				control := itinerary.ControlPoint{
					Radius:   v.Radius,
					MaxSpeed: v.MaxSpeed,
					ETA:      v.ETA,
					Point:    point,
					Name:     v.Name,
					Type:     v.Type,
				}
				controlPoints = append(controlPoints, control)
			}

			iti := itinerary.Itinerary{
				ID:             msg.ID,
				OrganizationID: msg.OrganizationID,
				PaymentID:      msg.PaymentID,
				Route:          itiPlatform.RouteID,
				Polyline:       itiPlatform.GetPolyline(),
				ControlPoints:  controlPoints,
			}

			fmt.Printf("iti: %v\n", iti)

			if ctx.Sender() != nil {
				ctx.Respond(&app.MsgItinirary{Data: iti})
			}

			return nil

		}(); err != nil {
			logs.LogError.Println(err)
			fmt.Printf("GetItinerary err: %s\n", err)
			return
		}
	}
}

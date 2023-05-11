package restclient

import (
	"fmt"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-driverconsole/internal/app"
)

func TestNewActor(t *testing.T) {
	type args struct {
		id           string
		user         string
		pass         string
		url          string
		keyUrl       string
		realm        string
		clientid     string
		clientsecret string
	}
	tests := []struct {
		name     string
		args     args
		messages []interface{}
	}{
		{
			name: "test1",
			args: args{
				id:           "NE-RC9-2957",
				user:         "device.sibus@sibus.com",
				pass:         "TuaeUs2DJpy",
				url:          "https://sibus.ambq.gov.co",
				keyUrl:       "https://sibus.ambq.gov.co/auth",
				realm:        "FLEET",
				clientid:     "devices-2",
				clientsecret: "fe889bd6-a28e-427a-8c2a-f5c73cd6f2d3",
			},
			messages: []interface{}{
				&app.MsgGetItinieary{
					ID:             "c28492bd-f114-416e-970d-dc8c8f9eab4e",
					OrganizationID: "11faeccb-148a-49ce-920b-33ac5807c2af",
					PaymentID:      114,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewActor(tt.args.id, tt.args.user, tt.args.pass, tt.args.url, tt.args.keyUrl, tt.args.realm, tt.args.clientid, tt.args.clientsecret)
			fmt.Println(got)
			if got == nil {
				t.Errorf("NewActor() = %v", got)
			}

			rootctx := actor.NewActorSystem().Root
			props := actor.PropsFromFunc(got.Receive)
			pid, err := rootctx.SpawnNamed(props, tt.name)
			if err != nil {
				t.Fatal(err)
			}

			for _, v := range tt.messages {
				rootctx.Request(pid, v)
				time.Sleep(time.Second * 10)
			}
		})
	}
}

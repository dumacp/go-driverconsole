package app

import (
	"reflect"
	"testing"

	"github.com/dumacp/go-schservices/api/services"
)

func BenchmarkUpdateService(b *testing.B) {
	prev := &services.ScheduleService{
		Id:    "uno",
		State: services.State_STARTED.String(),
		Route: &services.Route{
			Id:   130,
			Name: "route",
		},
	}
	current := &services.ScheduleService{
		Id:    "uno",
		State: services.State_ABORTED.String(),
	}
	for i := 0; i < b.N; i++ {
		UpdateService(prev, current)
	}
}

func BenchmarkUpdateServiceStable(b *testing.B) {
	prev := &services.ScheduleService{
		Id:    "uno",
		State: services.State_STARTED.String(),
		Route: &services.Route{
			Id:   130,
			Name: "route",
		},
	}
	current := &services.ScheduleService{
		Id:    "uno",
		State: services.State_ABORTED.String(),
	}
	for i := 0; i < b.N; i++ {
		UpdateServiceStable(prev, current)
	}
}

func TestUpdateService(t *testing.T) {
	type args struct {
		prev    *services.ScheduleService
		current *services.ScheduleService
	}
	tests := []struct {
		name string
		args args
		want *services.ScheduleService
	}{
		{
			name: "test1",
			args: args{
				prev: &services.ScheduleService{
					Id:    "uno",
					State: services.State_STARTED.String(),
					Route: &services.Route{
						Id:   130,
						Name: "route",
					},
				},
				current: &services.ScheduleService{
					Id:    "uno",
					State: services.State_ABORTED.String(),
				},
			},
			want: &services.ScheduleService{
				Id:    "uno",
				State: services.State_ABORTED.String(),
				Route: &services.Route{
					Id:   130,
					Name: "route",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateService(tt.args.prev, tt.args.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateServiceStable(t *testing.T) {
	type args struct {
		prev    *services.ScheduleService
		current *services.ScheduleService
	}
	tests := []struct {
		name string
		args args
		want *services.ScheduleService
	}{
		{
			name: "test1",
			args: args{
				prev: &services.ScheduleService{
					Id:    "uno",
					State: services.State_STARTED.String(),
					Route: &services.Route{
						Id:   130,
						Name: "route",
					},
				},
				current: &services.ScheduleService{
					Id:    "uno",
					State: services.State_ABORTED.String(),
				},
			},
			want: &services.ScheduleService{
				Id:    "uno",
				State: services.State_ABORTED.String(),
				Route: &services.Route{
					Id:   130,
					Name: "route",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateServiceStable(tt.args.prev, tt.args.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateService() = %v, want %v", got, tt.want)
			}
		})
	}
}

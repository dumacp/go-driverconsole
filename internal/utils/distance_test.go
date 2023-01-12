package utils

import (
	"testing"

	"github.com/golang/geo/s1"
)

func TestAngleToMeters(t *testing.T) {
	type args struct {
		a s1.Angle
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				a: 0.00003123,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AngleToMeters(tt.args.a); got != tt.want {
				t.Errorf("mts: %f", tt.args.a.Degrees()*111.139*1000)
				t.Errorf("AngleToMeters() = %v, want %v", got, tt.want)
			}
		})
	}
}

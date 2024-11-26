package main

import (
	"reflect"
	"testing"
)

func Test_getCandleDate(t *testing.T) {
	type args struct {
		code  string
		year  uint16
		month uint16
		day   uint16
		hour  uint16
	}
	tests := []struct {
		name    string
		args    args
		want    CandleResponse
		wantErr bool
	}{
		{
			name: "例題の動作担保",
			args: args{
				code:  "FTHD",
				year:  2021,
				month: 12,
				day:   22,
				hour:  10,
			},
			want: CandleResponse{
				Open:  3122,
				High:  3177,
				Low:   2865,
				Close: 2924,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCandleDate(tt.args.code, tt.args.year, tt.args.month, tt.args.day, tt.args.hour)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCandleDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCandleDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

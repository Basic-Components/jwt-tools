package jwtcentersdk

import (
	"testing"
	"time"
)

func TestAsymmetric_SignJSON(t *testing.T) {
	type args struct {
		jsonstringpayload string
		aud               string
		iss               string
	}
	tests := []struct {
		name string
		algo string
		args args
	}{
		{
			name: "测试",
			algo: "RS256",
			args: args{
				jsonstringpayload: `{"a":1}`,
				aud:               "1",
				iss:               "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remote, err := New("localhost:5000", tt.algo, time.Second)
			if err != nil {
				t.Errorf("remote init error = %v", err)
				return
			}
			got, err := remote.SignJSONString(tt.args.jsonstringpayload, tt.args.aud, tt.args.iss)
			if err != nil {
				t.Errorf("SignJSONString() error = %v", err)
				return
			}
			claims, err := remote.Verify(got)
			if err != nil {
				t.Errorf("Verify error = %v", err)
				return
			}
			if int(claims["a"].(float64)) != 1 {
				t.Errorf("want a = 1")
			}
		})
	}
}

func TestSymmetric_SignJSONString(t *testing.T) {
	type args struct {
		jsonstringpayload string
		aud               string
		iss               string
	}
	tests := []struct {
		name string
		algo string
		args args
	}{
		{
			name: "测试",
			algo: "HS256",
			args: args{
				jsonstringpayload: `{"a":1}`,
				aud:               "1",
				iss:               "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remote, err := New("localhost:5000", tt.algo, time.Second)
			if err != nil {
				t.Errorf("remote init error = %v", err)
				return
			}
			got, err := remote.SignJSONString(tt.args.jsonstringpayload, tt.args.aud, tt.args.iss)
			if err != nil {
				t.Errorf("remote.SignJSONString() error = %v", err)
				return
			}
			claims, err := remote.Verify(got)
			if err != nil {
				t.Errorf("remote.Verify error = %v", err)
				return
			}
			if int(claims["a"].(float64)) != 1 {
				t.Errorf("want a = 1")
			}
		})
	}
}

package jwtproxy

import (
	"testing"
)

func TestAsymmetric_SignJSON(t *testing.T) {
	type args struct {
		jsonstringpayload string
		aud               string
		iss               string
	}
	tests := []struct {
		name           string
		algo           string
		publishkeypath string
		privatekeypath string
		args           args
	}{
		{
			name:           "测试",
			algo:           "RS256",
			publishkeypath: "../autogen_rsa_pub.pem",
			privatekeypath: "../autogen_rsa.pem",
			args: args{
				jsonstringpayload: `{"a":1}`,
				aud:               "1",
				iss:               "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jwt = NewJwtProxy()
			err := jwt.InitAsymmetric(tt.algo, tt.privatekeypath, tt.publishkeypath)
			if err != nil {
				t.Errorf("AsymmetricFromPEMFile() error = %v", err)
				return
			}

			got, err := jwt.SignJSONString(tt.args.jsonstringpayload, tt.args.aud, tt.args.iss)
			if err != nil {
				t.Errorf("jwt.SignJSONString() error = %v", err)
				return
			}
			claims, err := jwt.Verify(got)
			if err != nil {
				t.Errorf("jwt.Verify error = %v", err)
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
		key  string
		args args
	}{
		{
			name: "测试",
			algo: "HS256",
			key:  "abcd",
			args: args{
				jsonstringpayload: `{"a":1}`,
				aud:               "1",
				iss:               "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jwt = NewJwtProxy()
			err := jwt.InitSymmetric(tt.algo, tt.key)
			if err != nil {
				t.Errorf("AsymmetricFromPEMFile() error = %v", err)
				return
			}

			got, err := jwt.SignJSONString(tt.args.jsonstringpayload, tt.args.aud, tt.args.iss)
			if err != nil {
				t.Errorf("jwt.SignJSONString() error = %v", err)
				return
			}
			claims, err := jwt.Verify(got)
			if err != nil {
				t.Errorf("jwt.Verify error = %v", err)
				return
			}
			if int(claims["a"].(float64)) != 1 {
				t.Errorf("want a = 1")
			}
		})
	}
}

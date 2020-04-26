package jwtsigner

import (
	"testing"

	"github.com/Basic-Components/jwttools/jwtverifier"
)

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
			signer, err := SymmetricNew(tt.algo, tt.key)
			if err != nil {
				t.Errorf("Symmetric.SignJSONString() error = %v", err)
				return
			}
			jwtverifier := jwtverifier.SymmetricNew(tt.key)
			got, err := signer.SignJSONString(tt.args.jsonstringpayload, tt.args.aud, tt.args.iss)
			if err != nil {
				t.Errorf("Symmetric.SignJSONString() error = %v", err)
				return
			}
			claims, err := jwtverifier.Verify(got)
			if int(claims["a"].(float64)) != 1 {
				t.Errorf("want a = 1")
			}
		})
	}
}

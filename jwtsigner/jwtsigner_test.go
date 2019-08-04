package jwtsigner

import (
	"fmt"
	"testing"
)

func TestSigner_Sign(t *testing.T) {
	type fields struct {
		method  string
		keyPath string
		claims  map[string]interface{}
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				method:  "RS256",
				keyPath: "../autogen_rsa.pem",
				claims:  map[string]interface{}{"my": 1}},
			args: args{
				data: []byte(`{"IP": "127.0.0.1", "name": "SKY"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signer, err := NewFromPath(tt.fields.method, tt.fields.keyPath, tt.fields.claims)
			if err != nil {
				t.Errorf("Signer NewFromPath() error = %v", err)
				return
			}
			got, err := signer.Sign(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signer.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("token:   %v\n", got)
		})
	}
}

package jwtclient

import (
	"testing"
)

func TestClient_Query(t *testing.T) {
	type fields struct {
		Address string
	}
	tests := []struct {
		name   string
		fields fields
		claims map[string]interface{}
	}{
		{
			name: "test1",
			fields: fields{
				Address: "localhost:5000",
			},
			claims: map[string]interface{}{"IP": "127.0.0.1", "name": "SKY"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.fields.Address)
			token, err := client.GetToken(tt.claims)
			if err != nil {
				t.Errorf("got token error: %v", err)
				return
			}
			gotclaims, err := client.VerifyToken(token)
			if err != nil {
				t.Errorf("got claims error: %v", err)
				return
			}
			if gotclaims["IP"] == tt.claims["IP"] && gotclaims["name"] == tt.claims["name"] {
				t.Logf("ok got claims: %v", gotclaims)
			}
		})
	}
}

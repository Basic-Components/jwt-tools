package jwtsigner

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func makeclaims(payload map[string]interface{}, aud string, iss string, exp int64) jwt.MapClaims {
	var claims jwt.MapClaims = payload
	now := time.Now().Unix()
	if aud != "" {
		claims["aud"] = aud
	}
	if iss != "" {
		claims["iss"] = iss
	}
	if exp > 0 {
		claims["exp"] = now + exp
	}
	return claims
}

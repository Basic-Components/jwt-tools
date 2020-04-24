package jwtsigner

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func makclaims(payload map[string]interface{}, aud string, iss string, exp int64) jwt.MapClaims {
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

//Signer 签名器接口
type Signer interface {

	// Sign 签名一个无过期的token
	Sign(payload map[string]interface{}, aud string, iss string) (string, error)

	// ExpSign 签名一个会过期的token
	ExpSign(payload map[string]interface{}, aud string, iss string, exp int64) (string, error)

	// SignJSON 为json签名一个无过期的token
	SignJSON(jsonpayload []byte, aud string, iss string) (string, error)

	// ExpSignJSON 为json签名一个无过期的token
	ExpSignJSON(jsonpayload []byte, aud string, iss string, exp int64) (string, error)

	// SignJSONString 为json字符串签名一个无过期的token
	SignJSONString(jsonstringpayload string, aud string, iss string) (string, error)

	// ExpSignJSONString 为json字符串签名一个会过期的token
	ExpSignJSONString(jsonstringpayload string, aud string, iss string, exp int64) (string, error)
}

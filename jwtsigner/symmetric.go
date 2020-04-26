package jwtsigner

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	uuid "github.com/satori/go.uuid"
)

// Symmetric 对称加密签名器
type Symmetric struct {
	key string
	alg jwt.SigningMethod
}

// SymmetricMethods 对称加密支持的算法范围
var SymmetricMethods = map[string]bool{
	"HS256": true,
	"HS384": true,
	"HS512": true,
}

// SymmetricNew 创建一个非对称加密签名器对象
func SymmetricNew(method string, key string) (*Symmetric, error) {
	_, ok := SymmetricMethods[method]
	if !ok {
		return nil, ErrUnexpectedAlgo
	}
	alg := jwt.GetSigningMethod(method)
	//var signer *Signer
	signer := &Symmetric{
		key: key,
		alg: alg}
	return signer, nil
}

func (signer *Symmetric) signany(claims jwt.MapClaims) (string, error) {
	now := time.Now().Unix()
	claims["jti"] = uuid.NewV4().String()
	claims["iat"] = now
	token := jwt.NewWithClaims(signer.alg, claims)
	out, err := token.SignedString([]byte(signer.key))
	if err != nil {
		return "", err
	}
	return out, nil
	//return "", ErrSignToken
}

// Sign 签名一个无过期的token
func (signer *Symmetric) Sign(payload map[string]interface{}, aud string, iss string) (string, error) {
	claims := makeclaims(payload, aud, iss, 0)
	return signer.signany(claims)
}

// ExpSign 签名一个会过期的token
func (signer *Symmetric) ExpSign(payload map[string]interface{}, aud string, iss string, exp int64) (string, error) {
	claims := makeclaims(payload, aud, iss, exp)
	return signer.signany(claims)
}

// SignJSON 为json签名一个无过期的token
func (signer *Symmetric) SignJSON(jsonpayload []byte, aud string, iss string) (string, error) {
	payload := map[string]interface{}{}
	err := json.Unmarshal(jsonpayload, &payload)
	if err != nil {
		return "", err //ErrParseClaimsToJSON
	}
	return signer.Sign(payload, aud, iss)
}

// ExpSignJSON 为json签名一个无过期的token
func (signer *Symmetric) ExpSignJSON(jsonpayload []byte, aud string, iss string, exp int64) (string, error) {
	payload := map[string]interface{}{}
	err := json.Unmarshal(jsonpayload, &payload)
	if err != nil {
		return "", err //ErrParseClaimsToJSON
	}
	return signer.ExpSign(payload, aud, iss, exp)
}

// SignJSONString 为json字符串签名一个无过期的token
func (signer *Symmetric) SignJSONString(jsonstringpayload string, aud string, iss string) (string, error) {
	jsonpayload := []byte(jsonstringpayload)
	return signer.SignJSON(jsonpayload, aud, iss)
}

// ExpSignJSONString 为json字符串签名一个会过期的token
func (signer *Symmetric) ExpSignJSONString(jsonstringpayload string, aud string, iss string, exp int64) (string, error) {
	jsonpayload := []byte(jsonstringpayload)
	return signer.ExpSignJSON(jsonpayload, aud, iss, exp)
}

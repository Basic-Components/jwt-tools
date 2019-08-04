package jwtsigner

import (
	"encoding/json"
	"time"

	errs "github.com/Basic-Components/jwtrpc/errs"
	utils "github.com/Basic-Components/jwtrpc/utils"

	jwt "github.com/dgrijalva/jwt-go"

	uuid "github.com/satori/go.uuid"
)

// Signer 签名器
type Signer struct {
	key           interface{}
	alg           jwt.SigningMethod
	defaultClaims map[string]interface{}
}

// New 创建一个Signer对象
func New(method string, key interface{}, claims map[string]interface{}) *Signer {
	alg := jwt.GetSigningMethod(method)
	//var signer *Signer
	signer := &Signer{
		key:           key,
		alg:           alg,
		defaultClaims: claims}
	return signer
}

// NewFromPath 从路径上读取密钥文件创建一个Signer对象
func NewFromPath(method string, keyPath string, claims map[string]interface{}) (*Signer, error) {
	keybytes, err := utils.LoadData(keyPath)

	if err != nil {
		return nil, errs.LoadKeyError
	}
	if utils.IsEs(method) {
		key, err := jwt.ParseECPrivateKeyFromPEM(keybytes)
		if err != nil {
			return nil, err
		}
		return New(method, key, claims), nil
	} else if utils.IsRs(method) {
		key, err := jwt.ParseRSAPrivateKeyFromPEM(keybytes)
		if err != nil {
			return nil, err
		}
		return New(method, key, claims), nil
	} else {
		return nil, errs.SignerTypeError
	}
}

// Sign 用Signer对象签名
func (signer *Signer) Sign(data []byte) (string, error) {
	var claims jwt.MapClaims
	if err := json.Unmarshal(data, &claims); err != nil {
		return "", errs.ParseClaimsToJsonError
	} else {
		for key, value := range signer.defaultClaims {
			claims[key] = value
		}
		claims["iat"] = int(time.Now().Unix())
		claims["jti"] = uuid.NewV4().String()
		token := jwt.NewWithClaims(signer.alg, claims)
		if out, err := token.SignedString(signer.key); err == nil {
			return out, nil
		} else {
			return "", errs.SignTokenError
		}
	}
}

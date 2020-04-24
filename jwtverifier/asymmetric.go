package jwtverifier

import (
	"regexp"

	errs "github.com/Basic-Components/jwtrpc/errs"
	utils "github.com/Basic-Components/jwtrpc/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

// Asymmetric 非对称加密jwt的验证器
type Asymmetric struct {
	key interface{}
	alg jwt.SigningMethod
}

// AsymmetricMethods 非对称加密支持的算法范围
var AsymmetricMethods = map[string]bool{
	"RS256": true,
	"RS384": true,
	"RS512": true,
	"ES256": true,
	"ES384": true,
	"ES512": true,
}

// asymmetricNew 创建一个非对称加密jwt的验证器对象
func asymmetricNew(method string, key interface{}) (*Asymmetric, error) {
	_, ok := AsymmetricMethods[method]
	if !ok {
		return nil, ErrUnexpectedAlgo
	}
	alg := jwt.GetSigningMethod(method)
	var verifier *Asymmetric
	verifier = &Asymmetric{
		key: key,
		alg: alg,
	}
	return verifier, nil
}

// AsymmetricFromPEM 从路径上读取密钥文件创建一个Verifier对象
func AsymmetricFromPEM(method string, keyPath string) (*Asymmetric, error) {
	keybytes, err := utils.LoadData(keyPath)
	if err != nil {
		return nil, errs.LoadKeyError
	}
	if utils.IsEs(method) {
		key, err := jwt.ParseECPublicKeyFromPEM(keybytes)
		if err != nil {
			return nil, err
		}
		return asymmetricNew(method, key)
	} else if utils.IsRs(method) {
		key, err := jwt.ParseRSAPublicKeyFromPEM(keybytes)
		if err != nil {
			return nil, err
		}
		return asymmetricNew(method, key)
	} else {
		return nil, ErrUnexpectedAlgo
	}
}

// Verify 用Verifier对象验签
func (verifier *Asymmetric) Verify(tokenstring string) (map[string]interface{}, error) {
	tokenBytes := []byte(tokenstring)
	tokData := regexp.MustCompile(`\s*$`).ReplaceAll(tokenBytes, []byte{})
	token, err := jwt.Parse(
		string(tokData),
		func(t *jwt.Token) (interface{}, error) {
			return verifier.key, nil
		})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		var payload map[string]interface{} = claims
		return payload, nil
	}
	return nil, ErrVerifyToken
}

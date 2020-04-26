package jwtverifier

import (
	"regexp"

	utils "github.com/Basic-Components/jwttools/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

// Asymmetric 非对称加密jwt的验证器
type Asymmetric struct {
	key interface{}
	alg jwt.SigningMethod
}

// asymmetricNew 创建一个非对称加密jwt的验证器对象
func asymmetricNew(method string, key interface{}) (*Asymmetric, error) {
	_, ok := utils.AsymmetricMethods[method]
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

// AsymmetricFromPEM 使用PEM编码的密钥字节串创建一个非对称加密的jwt验证器对象
func AsymmetricFromPEM(method string, keybytes []byte) (*Asymmetric, error) {
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

// AsymmetricFromPEMFile 从路径上读取密钥文件创建一个Verifier对象
func AsymmetricFromPEMFile(method string, keyPath string) (*Asymmetric, error) {
	keybytes, err := utils.LoadData(keyPath)
	if err != nil {
		return nil, ErrLoadPublicKey
	}
	return AsymmetricFromPEM(method, keybytes)
}

// Verify 用Verifier对象验签
func (verifier *Asymmetric) Verify(tokenstring string) (map[string]interface{}, error) {
	tokenBytes := []byte(tokenstring)
	tokData := regexp.MustCompile(`\s*$`).ReplaceAll(tokenBytes, []byte{})
	token, err := jwt.Parse(
		string(tokData),
		func(t *jwt.Token) (interface{}, error) {
			// if utils.IsEs(verifier.alg.Alg()) {
			// 	return verifier.key.(*ecdsa.PublicKey), nil
			// }
			// if utils.IsRs(verifier.alg.Alg()) {
			// 	return verifier.key.(*rsa.PublicKey), nil
			// }
			// return nil, ErrUnexpectedAlgo
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

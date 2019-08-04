package jwtverifier

import (
	"regexp"

	errs "github.com/Basic-Components/jwtrpc/errs"
	utils "github.com/Basic-Components/jwtrpc/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

// Verifier jwt的验证器
type Verifier struct {
	key interface{}
	alg jwt.SigningMethod
}

// New 创建一个Verifier对象
func New(method string, key interface{}) *Verifier {
	alg := jwt.GetSigningMethod(method)
	var verifier *Verifier
	verifier = &Verifier{
		key: key,
		alg: alg,
	}
	return verifier
}

// NewFromPath 从路径上读取密钥文件创建一个Verifier对象
func NewFromPath(method string, keyPath string) (*Verifier, error) {
	keybytes, err := utils.LoadData(keyPath)
	if err != nil {
		return nil, errs.LoadKeyError
	}
	if utils.IsEs(method) {
		key, err := jwt.ParseECPublicKeyFromPEM(keybytes)
		if err != nil {
			return nil, err
		}
		return New(method, key), nil
	} else if utils.IsRs(method) {
		key, err := jwt.ParseRSAPublicKeyFromPEM(keybytes)
		if err != nil {
			return nil, err
		}
		return New(method, key), nil
	} else {
		return nil, errs.SignerTypeError
	}
}

// Verify 用Verifier对象验签
func (verifier *Verifier) Verify(tokData []byte) (jwt.MapClaims, error) {
	tokData = regexp.MustCompile(`\s*$`).ReplaceAll(tokData, []byte{})
	token, err := jwt.Parse(string(tokData), func(t *jwt.Token) (interface{}, error) {
		return verifier.key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return claims, errs.VerifyTokenError
}

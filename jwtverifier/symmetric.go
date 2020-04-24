package jwtverifier

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// Symmetric jwt的对称的加密验证器
type Symmetric struct {
	key string
}

// SymmetricNew 创建一个对称的加密验证器对象
func SymmetricNew(key string) *Symmetric {
	verifier := &Symmetric{
		key: key}
	return verifier
}

// Verify 用Verifier对象验签
func (verifier *Symmetric) Verify(tokenstring string) (map[string]interface{}, error) {
	token, err := jwt.Parse(
		tokenstring,
		func(t *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			return []byte(verifier.key), nil
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

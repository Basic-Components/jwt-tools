package jwtverifier

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//Verifier 验证器接口
type Verifier interface {

	// Verify 用Verifier对象验签
	Verify(tokenstring string) (map[string]interface{}, error)
}

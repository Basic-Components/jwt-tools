package interfaces

import (
	"crypto"
	"io"

	"github.com/Basic-Components/jwttools/options"
)

//CanSign 签名器接口
type CanSign interface {
	// Sign 签名一个token
	Sign(payload []byte, opts ...options.SignOption) (string, error)
}


//CanVerify 验证器接口
type CanVerify interface {
	// Verify 用Verifier对象验签
	Verify(tokenstring string) (map[string]interface{}, error)
}

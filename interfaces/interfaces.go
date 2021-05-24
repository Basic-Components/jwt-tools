package interfaces

import (
	declare "github.com/Basic-Components/jwttools/jwtrpcdeclare"
	"github.com/Basic-Components/jwttools/options"
)

//CanSign 签名器接口
type CanSign interface {
	// Sign 签名一个token
	Sign(payload []byte, opts ...options.SignOption) (string, error)
	//Meta 查看签名器元信息
	Meta() *declare.MetaResponse
}

//CanVerify 验证器接口
type CanVerify interface {
	// Verify 用Verifier对象验签
	Verify(tokenstring string) (map[string]interface{}, error)
}

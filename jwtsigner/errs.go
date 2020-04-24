package jwtsigner

import (
	"errors"
)

// ErrUnexpectedAlgo 算法类型错误
var ErrUnexpectedAlgo = errors.New("unknown algo type key")

// ErrLoadPrivateKey 私钥无法阅读
var ErrLoadPrivateKey = errors.New("couldn't read key")

// ErrParseClaimsToJSON 无法加载JSON
var ErrParseClaimsToJSON = errors.New("Couldn't parse claims JSON")

// ErrSignToken 签名错误
var ErrSignToken = errors.New("Error signing token")

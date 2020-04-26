package errs

import (
	"errors"
)

var ErrLoadKey = errors.New("couldn't read key")

var TokenInvalidError error = errors.New("Token is invalid")
var VerifyTokenError error = errors.New("Verify Token error")

// ErrConfigParams 配置参数错误
var ErrConfigParams = errors.New("config params error")

// ErrAlgoType 算法不支持
var ErrAlgoType = errors.New("unknown algo type key")

// ErrExpOutOfRange 过期时间超出范围
var ErrExpOutOfRange = errors.New("exp is out of range")

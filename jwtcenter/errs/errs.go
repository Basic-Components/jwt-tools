package errs

import (
	"errors"
)

var ErrAlgoType = errors.New("unknown algo type key")
var ErrLoadKey = errors.New("couldn't read key")

var TokenInvalidError error = errors.New("Token is invalid")
var VerifyTokenError error = errors.New("Verify Token error")

// ErrConfigParams 配置参数错误
var ErrConfigParams = errors.New("config params error")

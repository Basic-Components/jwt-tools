package errs

import (
	"errors"
)

// ErrAlgoType 算法类型不支持
var ErrAlgoType = errors.New("unknown algo type key")

// ErrProxyNotInited 数据库代理未初始化错误
var ErrProxyNotInited = errors.New("proxy not inited yet")

// ErrProxyAlreadyInited 代理已经被初始化过了
var ErrProxyAlreadyInited = errors.New("proxy already inited yet")

package errs

import (
	"errors"
)

// ErrAlgoType 算法类型不支持
var ErrAlgoType = errors.New("unknown algo type key")

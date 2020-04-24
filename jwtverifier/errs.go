package jwtverifier

import (
	"errors"
)

// ErrUnexpectedAlgo 算法不符
var ErrUnexpectedAlgo = errors.New("unexpected algo type")

// ErrVerifyToken 校验token错误
var ErrVerifyToken = errors.New("Verify Token error")

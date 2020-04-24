package jwtverifier

import (
	"errors"
)

// ErrUnexpectedAlgo 算法不符
var ErrUnexpectedAlgo = errors.New("unexpected algo type")

// ErrVerifyToken 校验token错误
var ErrVerifyToken = errors.New("Verify Token error")

// ErrLoadPublicKey 公钥无法阅读
var ErrLoadPublicKey = errors.New("couldn't read public key")

package jwtproxy

import (
	"errors"

	"github.com/Basic-Components/jwttools/jwtsigner"
	"github.com/Basic-Components/jwttools/jwtverifier"
	"github.com/Basic-Components/jwttools/utils"
)

// ErrProxyNotInited 代理未初始化错误
var ErrProxyNotInited = errors.New("db proxy not inited yet")

// ErrProxyAlreadyInited 代理已经被初始化过了
var ErrProxyAlreadyInited = errors.New("proxy already inited yet")

// DBProxyCallback 数据库操作的回调函数
type jwtProxyCallback func(signer utils.Signer, verifier utils.Verifier) error

// jwtProxy 数据库客户端的代理
type jwtProxy struct {
	Ok        bool
	signer    utils.Signer
	verifier  utils.Verifier
	callBacks []jwtProxyCallback
}

// NewJwtProxy 创建一个新的数据库客户端代理
func NewJwtProxy() *jwtProxy {
	proxy := new(jwtProxy)
	proxy.Ok = false
	return proxy
}

// Jwt 默认的etcd代理对象
var Jwt = NewJwtProxy()

func (proxy *jwtProxy) runCallback() {
	for _, cb := range proxy.callBacks {
		cb(proxy.signer, proxy.verifier)
	}
}

// InitSymmetric 初始化对称加密的签名器和验证器
func (proxy *jwtProxy) InitSymmetric(method string, key string) error {
	if proxy.Ok {
		return ErrProxyAlreadyInited
	}
	signer, err := jwtsigner.SymmetricNew(method, key)
	if err != nil {
		return err
	}
	verifier := jwtverifier.SymmetricNew(key)
	proxy.signer = signer
	proxy.verifier = verifier
	proxy.runCallback()
	proxy.Ok = true
	return nil
}

// InitFromURL 初始化非对称加密的签名器和验证器
func (proxy *jwtProxy) InitAsymmetric(method string, privatekeyPath string, publickeyPath string) error {
	if proxy.Ok {
		return ErrProxyAlreadyInited
	}
	signer, err := jwtsigner.AsymmetricFromPEMFile(method, privatekeyPath)
	if err != nil {
		return err
	}
	verifier, err := jwtverifier.AsymmetricFromPEMFile(method, publickeyPath)
	if err != nil {
		return err
	}
	proxy.signer = signer
	proxy.verifier = verifier
	proxy.runCallback()
	proxy.Ok = true
	return nil
}

// Regist 注册回调函数,在init执行后执行回调函数
func (proxy *jwtProxy) Regist(cb jwtProxyCallback) {
	proxy.callBacks = append(proxy.callBacks, cb)
}

// SignJSON 为json签名一个无过期的token
func (proxy *jwtProxy) SignJSON(jsonpayload []byte, aud string, iss string) (string, error) {
	if !proxy.Ok {
		return "", ErrProxyNotInited
	}
	return proxy.signer.SignJSON(jsonpayload, aud, iss)

}

// ExpSignJSON 为json签名一个会过期的token
func (proxy *jwtProxy) ExpSignJSON(jsonpayload []byte, aud string, iss string, exp int64) (string, error) {
	if !proxy.Ok {
		return "", ErrProxyNotInited
	}
	return proxy.signer.ExpSignJSON(jsonpayload, aud, iss, exp)
}

// SignJSONString 为json字符串签名一个无过期的token
func (proxy *jwtProxy) SignJSONString(jsonstringpayload string, aud string, iss string) (string, error) {
	if !proxy.Ok {
		return "", ErrProxyNotInited
	}
	return proxy.signer.SignJSONString(jsonstringpayload, aud, iss)
}

// ExpSignJSONString 为json字符串签名一个会过期的token
func (proxy *jwtProxy) ExpSignJSONString(jsonstringpayload string, aud string, iss string, exp int64) (string, error) {
	if !proxy.Ok {
		return "", ErrProxyNotInited
	}
	return proxy.signer.ExpSignJSONString(jsonstringpayload, aud, iss, exp)
}

// Sign 签名一个无过期的token
func (proxy *jwtProxy) Sign(payload map[string]interface{}, aud string, iss string) (string, error) {
	if !proxy.Ok {
		return "", ErrProxyNotInited
	}
	return proxy.signer.Sign(payload, aud, iss)
}

// ExpSign 签名一个会过期的token
func (proxy *jwtProxy) ExpSign(payload map[string]interface{}, aud string, iss string, exp int64) (string, error) {
	if !proxy.Ok {
		return "", ErrProxyNotInited
	}
	return proxy.signer.ExpSign(payload, aud, iss, exp)
}

// Verify 验证签名
func (proxy *jwtProxy) Verify(token string) (map[string]interface{}, error) {
	if !proxy.Ok {
		return nil, ErrProxyNotInited
	}
	return proxy.verifier.Verify(token)
}

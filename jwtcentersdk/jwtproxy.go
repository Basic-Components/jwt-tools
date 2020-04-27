package jwtcentersdk

import (
	"time"

	errs "github.com/Basic-Components/jwttools/jwtcentersdk/errs"
)

// jwtProxyCallback jwt代理注册后执行的回调函数

type jwtProxyCallback func(rc *RemoteCenter) error

//初始化代理的参数
type JwtOptions struct {
	Method  string
	Address string
	Timeout time.Duration
}

// jwtProxy 数据库客户端的代理
type jwtProxy struct {
	Ok        bool
	rc        *RemoteCenter
	callBacks []jwtProxyCallback
}

// NewJwtProxy 创建一个新的数据库客户端代理
func NewJwtProxy() *jwtProxy {
	proxy := new(jwtProxy)
	proxy.Ok = false
	return proxy
}

// Jwt 默认的jwtcentersdk代理对象
var Jwt = NewJwtProxy()

// Close 关闭连接
func (proxy *jwtProxy) Close() {
	if proxy.Ok {
		proxy.rc.Close()
	}
	proxy.Ok = false
}

func (proxy *jwtProxy) runCallback() {
	for _, cb := range proxy.callBacks {
		//go cb(proxy.Cli)
		cb(proxy.rc)
	}
}

// Init 使用配置给代理赋值客户端实例
func (proxy *jwtProxy) Init(address string, algo string, timeout time.Duration) error {
	if proxy.Ok {
		return errs.ErrProxyAlreadyInited
	}
	rc, err := New(address, algo, timeout)
	if err != nil {
		return err
	}
	proxy.rc = rc
	proxy.runCallback()
	proxy.Ok = true
	return nil
}

// InitFromURL 使用配置给代理赋值客户端实例
func (proxy *jwtProxy) InitWithLocalBalance(addresses []string, algo string, timeout time.Duration) error {
	if proxy.Ok {
		return errs.ErrProxyAlreadyInited
	}
	rc, err := NewWithLocalBalance(addresses, algo, timeout)
	if err != nil {
		return err
	}
	proxy.rc = rc
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
		return "", errs.ErrProxyNotInited
	}
	return proxy.rc.SignJSON(jsonpayload, aud, iss)

}

// ExpSignJSON 为json签名一个会过期的token
func (proxy *jwtProxy) ExpSignJSON(jsonpayload []byte, aud string, iss string, exp int64) (string, error) {
	if !proxy.Ok {
		return "", errs.ErrProxyNotInited
	}
	return proxy.rc.ExpSignJSON(jsonpayload, aud, iss, exp)
}

// SignJSONString 为json字符串签名一个无过期的token
func (proxy *jwtProxy) SignJSONString(jsonstringpayload string, aud string, iss string) (string, error) {
	if !proxy.Ok {
		return "", errs.ErrProxyNotInited
	}
	return proxy.rc.SignJSONString(jsonstringpayload, aud, iss)
}

// ExpSignJSONString 为json字符串签名一个会过期的token
func (proxy *jwtProxy) ExpSignJSONString(jsonstringpayload string, aud string, iss string, exp int64) (string, error) {
	if !proxy.Ok {
		return "", errs.ErrProxyNotInited
	}
	return proxy.rc.ExpSignJSONString(jsonstringpayload, aud, iss, exp)
}

// Sign 签名一个无过期的token
func (proxy *jwtProxy) Sign(payload map[string]interface{}, aud string, iss string) (string, error) {
	if !proxy.Ok {
		return "", errs.ErrProxyNotInited
	}
	return proxy.rc.Sign(payload, aud, iss)
}

// ExpSign 签名一个会过期的token
func (proxy *jwtProxy) ExpSign(payload map[string]interface{}, aud string, iss string, exp int64) (string, error) {
	if !proxy.Ok {
		return "", errs.ErrProxyNotInited
	}
	return proxy.rc.ExpSign(payload, aud, iss, exp)
}

// Verify 验证签名
func (proxy *jwtProxy) Verify(token string) (map[string]interface{}, error) {
	if !proxy.Ok {
		return nil, errs.ErrProxyNotInited
	}
	return proxy.rc.Verify(token)
}

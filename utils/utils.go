package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//Signer 签名器接口
type Signer interface {

	// Sign 签名一个无过期的token
	Sign(payload map[string]interface{}, aud string, iss string) (string, error)

	// ExpSign 签名一个会过期的token
	ExpSign(payload map[string]interface{}, aud string, iss string, exp int64) (string, error)

	// SignJSON 为json签名一个无过期的token
	SignJSON(jsonpayload []byte, aud string, iss string) (string, error)

	// ExpSignJSON 为json签名一个会过期的token
	ExpSignJSON(jsonpayload []byte, aud string, iss string, exp int64) (string, error)

	// SignJSONString 为json字符串签名一个无过期的token
	SignJSONString(jsonstringpayload string, aud string, iss string) (string, error)

	// ExpSignJSONString 为json字符串签名一个会过期的token
	ExpSignJSONString(jsonstringpayload string, aud string, iss string, exp int64) (string, error)
}

//Verifier 验证器接口
type Verifier interface {

	// Verify 用Verifier对象验签
	Verify(tokenstring string) (map[string]interface{}, error)
}

// AsymmetricMethods 非对称加密支持的算法范围
var AsymmetricMethods = map[string]bool{
	"RS256": true,
	"RS384": true,
	"RS512": true,
	"ES256": true,
	"ES384": true,
	"ES512": true,
}

// SymmetricMethods 对称加密支持的算法范围
var SymmetricMethods = map[string]bool{
	"HS256": true,
	"HS384": true,
	"HS512": true,
}

// AsymmetricMethods 非对称加密支持的算法范围
var CenterSupportedMethods = map[string]bool{
	"RS256": true,
	"HS256": true,
}

// IsEs 判断文件是不是ES方法加密
func IsEs(method string) bool {
	return strings.HasPrefix(method, "ES")
}

// IsRs 判断文件是不是RS方法加密
func IsRs(method string) bool {
	return strings.HasPrefix(method, "RS") || strings.HasPrefix(method, "PS")
}

// LoadData 读取并加载文件数据
func LoadData(p string) ([]byte, error) {
	if p == "" {
		return nil, fmt.Errorf("No path specified")
	}

	var rdr io.Reader
	if p == "-" {
		rdr = os.Stdin
	} else if p == "+" {
		return []byte("{}"), nil
	} else {
		if f, err := os.Open(p); err == nil {
			rdr = f
			defer f.Close()
		} else {
			return nil, err
		}
	}
	return ioutil.ReadAll(rdr)
}

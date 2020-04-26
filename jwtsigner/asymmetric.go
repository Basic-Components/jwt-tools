package jwtsigner

import (
	"crypto"
	"io"
	"time"

	utils "github.com/Basic-Components/jwttools/utils"
	jwt "github.com/dgrijalva/jwt-go"

	uuid "github.com/satori/go.uuid"
)

// PrivateKey 非对称加密的私钥
type PrivateKey interface {
	Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error)
}

// Asymmetric 非对称加密签名器
type Asymmetric struct {
	key PrivateKey
	alg jwt.SigningMethod
}

// AsymmetricNew 创建一个非对称加密签名器对象
func AsymmetricNew(method string, key PrivateKey) (*Asymmetric, error) {
	_, ok := utils.AsymmetricMethods[method]
	if !ok {
		return nil, ErrUnexpectedAlgo
	}
	alg := jwt.GetSigningMethod(method)
	//var signer *Signer
	signer := &Asymmetric{
		key: key,
		alg: alg}
	return signer, nil
}

// AsymmetricFromPEM 使用PEM编码的密钥字节串创建一个非对称加密签名器对象
func AsymmetricFromPEM(method string, keybytes []byte) (*Asymmetric, error) {
	if utils.IsEs(method) {
		key, err := jwt.ParseECPrivateKeyFromPEM(keybytes)
		if err != nil {
			return nil, err
		}
		return AsymmetricNew(method, key)
	} else if utils.IsRs(method) {
		key, err := jwt.ParseRSAPrivateKeyFromPEM(keybytes)
		if err != nil {
			return nil, err
		}
		return AsymmetricNew(method, key)
	} else {
		return nil, ErrUnexpectedAlgo
	}
}

// AsymmetricFromPEMFile 从路径上读取PEM密钥文件创建一个非对称加密签名器对象
func AsymmetricFromPEMFile(method string, keyPath string) (*Asymmetric, error) {
	keybytes, err := utils.LoadData(keyPath)
	if err != nil {
		return nil, ErrLoadPrivateKey
	}
	return AsymmetricFromPEM(method, keybytes)
}

func (signer *Asymmetric) signany(claims jwt.MapClaims) (string, error) {
	now := time.Now().Unix()
	claims["jti"] = uuid.NewV4().String()
	claims["iat"] = now
	token := jwt.NewWithClaims(signer.alg, claims)
	out, err := token.SignedString(signer.key)
	if err != nil {
		return "", err
		//return "",
	}
	return out, nil
}

// Sign 签名一个无过期的token
func (signer *Asymmetric) Sign(payload map[string]interface{}, aud string, iss string) (string, error) {
	claims := makeclaims(payload, aud, iss, 0)
	return signer.signany(claims)
}

// ExpSign 签名一个会过期的token
func (signer *Asymmetric) ExpSign(payload map[string]interface{}, aud string, iss string, exp int64) (string, error) {
	claims := makeclaims(payload, aud, iss, exp)
	return signer.signany(claims)
}

// SignJSON 为json签名一个无过期的token
func (signer *Asymmetric) SignJSON(jsonpayload []byte, aud string, iss string) (string, error) {
	payload := map[string]interface{}{}
	err := json.Unmarshal(jsonpayload, &payload)
	if err != nil {
		return "", err
	}
	return signer.Sign(payload, aud, iss)
}

// ExpSignJSON 为json签名一个无过期的token
func (signer *Asymmetric) ExpSignJSON(jsonpayload []byte, aud string, iss string, exp int64) (string, error) {
	payload := map[string]interface{}{}
	err := json.Unmarshal(jsonpayload, &payload)
	if err != nil {
		return "", err
	}
	return signer.ExpSign(payload, aud, iss, exp)
}

// SignJSONString 为json字符串签名一个无过期的token
func (signer *Asymmetric) SignJSONString(jsonstringpayload string, aud string, iss string) (string, error) {
	jsonpayload := []byte(jsonstringpayload)
	return signer.SignJSON(jsonpayload, aud, iss)
}

// ExpSignJSONString 为json字符串签名一个会过期的token
func (signer *Asymmetric) ExpSignJSONString(jsonstringpayload string, aud string, iss string, exp int64) (string, error) {
	jsonpayload := []byte(jsonstringpayload)
	return signer.ExpSignJSON(jsonpayload, aud, iss, exp)
}

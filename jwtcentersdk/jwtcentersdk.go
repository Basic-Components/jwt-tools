package jwtcentersdk

import (
	"context"
	"errors"
	"time"

	errs "github.com/Basic-Components/jwttools/jwtcenter/errs"
	pb "github.com/Basic-Components/jwttools/jwtcenter/jwtrpcdeclare"
	utils "github.com/Basic-Components/jwttools/utils"
	jsoniter "github.com/json-iterator/go"
	grpc "google.golang.org/grpc"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// RemoteCenter jwt的客户端类型
type RemoteCenter struct {
	Algo    pb.Algo
	Address string
	Timeout time.Duration
}

// New 创建客户端对象
func New(address string, algo string, timeout time.Duration) (*RemoteCenter, error) {
	rc := new(RemoteCenter)
	rc.Address = address
	rc.Timeout = timeout
	_, ok := utils.CenterSupportedMethods[algo]
	if !ok {
		return rc, errs.ErrAlgoType
	}
	switch algo {
	case "RS256":
		rc.Algo = pb.Algo_RS256
	default:
		rc.Algo = pb.Algo_HS256
	}
	return rc, nil
}

// SignJSON 为json签名一个无过期的token
func (client *RemoteCenter) SignJSON(jsonpayload []byte, aud string, iss string) (string, error) {
	conn, err := grpc.Dial(client.Address, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewJwtServiceClient(conn)
	// 设置请求上下文的过期时间
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()
	rs, err := c.SignJSON(ctx, &pb.SignJSONRequest{
		Algo:    client.Algo,
		Payload: jsonpayload,
		Aud:     aud,
		Iss:     iss,
	})
	if err != nil {
		return "", err
	}
	if rs.Status.Status == pb.StatusData_ERROR {
		return "", errors.New(rs.Status.Msg)
	}
	return rs.Token, nil
}

// ExpSignJSON 为json签名一个会过期的token
func (client *RemoteCenter) ExpSignJSON(jsonpayload []byte, aud string, iss string, exp int64) (string, error) {
	conn, err := grpc.Dial(client.Address, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewJwtServiceClient(conn)
	// 设置请求上下文的过期时间
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()
	rs, err := c.SignJSON(ctx, &pb.SignJSONRequest{
		Algo:    client.Algo,
		Payload: jsonpayload,
		Aud:     aud,
		Iss:     iss,
		Exp:     exp,
	})
	if err != nil {
		return "", err
	}
	if rs.Status.Status == pb.StatusData_ERROR {
		return "", errors.New(rs.Status.Msg)
	}
	return rs.Token, nil
}

// SignJSONString 为json字符串签名一个无过期的token
func (client *RemoteCenter) SignJSONString(jsonstringpayload string, aud string, iss string) (string, error) {
	jsonpayload := []byte(jsonstringpayload)
	return client.SignJSON(jsonpayload, aud, iss)
}

// ExpSignJSONString 为json字符串签名一个会过期的token
func (client *RemoteCenter) ExpSignJSONString(jsonstringpayload string, aud string, iss string, exp int64) (string, error) {
	jsonpayload := []byte(jsonstringpayload)
	return client.ExpSignJSON(jsonpayload, aud, iss, exp)
}

// Sign 签名一个无过期的token
func (client *RemoteCenter) Sign(payload map[string]interface{}, aud string, iss string) (string, error) {
	jsonpayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return client.SignJSON(jsonpayload, aud, iss)
}

// ExpSign 签名一个会过期的token
func (client *RemoteCenter) ExpSign(payload map[string]interface{}, aud string, iss string, exp int64) (string, error) {
	jsonpayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return client.ExpSignJSON(jsonpayload, aud, iss, exp)
}

// Verify 验证签名
func (client *RemoteCenter) Verify(token string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	conn, err := grpc.Dial(client.Address, grpc.WithInsecure())
	if err != nil {
		return result, err
	}
	defer conn.Close()
	c := pb.NewJwtServiceClient(conn)
	// 设置请求上下文的过期时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rs, err := c.VerifyJSON(ctx, &pb.VerifyRequest{Algo: client.Algo, Token: token})
	if err != nil {
		return result, err
	}
	if rs.Status.Status == pb.StatusData_ERROR {
		return result, errors.New(rs.Status.Msg)
	}
	err = json.Unmarshal([]byte(rs.Payload), &result)
	return result, err
}

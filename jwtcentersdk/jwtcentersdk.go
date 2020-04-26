package jwtcentersdk

import (
	"context"
	"errors"
	"time"

	pb "github.com/Basic-Components/jwttools/jwtcenter/jwtrpcdeclare"
	jsoniter "github.com/json-iterator/go"
	grpc "google.golang.org/grpc"
)

// Client jwt的客户端类型
type RemoteCenter struct {
	Algo    pb.Algo
	Address string
	Timeout int
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// New 创建客户端对象
func New(address string, algo string, timeout int) (*RemoteCenter, error) {
	rc := new(RemoteCenter)
	rc.Address = address
	rc.Timeout = timeout
	_, ok := CenterSupportedMethods[algo]
	if !ok {
		return rc, ErrAlgoType
	}
	switch algo {
	case "RS256":
		rc.Algo = pb.Algo_RS256
	default:
		rc.Algo = pb.Algo_HS256
	}
	return rc, nil
}

// 为json签名一个不会过期的token
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
		payload: jsonpayload,
		aud:     aud,
		iss:     iss,
	})
	if err != nil {
		return "", err
	}
	if rs.Status.Status == pb.StatusData_ERROR {
		return "", errors.New(rs.Status.Msg)
	}
	return rs.Token, nil
}

// 为json签名一个会过期的token
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
		payload: jsonpayload,
		aud:     aud,
		iss:     iss,
		exp:     exp,
	})
	if err != nil {
		return "", err
	}
	if rs.Status.Status == pb.StatusData_ERROR {
		return "", errors.New(rs.Status.Msg)
	}
	return rs.Token, nil
}

func (client *RemoteCenter) SignJSONString(jsonstringpayload string, aud string, iss string) (string, error) {
	jsonpayload := []byte(jsonstringpayload)
	return client.SignJSON(jsonpayload, aud, iss)
}

// ExpSignJSONString 为json字符串签名一个会过期的token
func (client *RemoteCenter) ExpSignJSONString(jsonstringpayload string, aud string, iss string, exp int64) (string, error) {
	jsonpayload := []byte(jsonstringpayload)
	return client.ExpSignJSON(jsonpayload, aud, iss, exp)
}

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
func (client *Client) Verify(token string) (map[string]interface{}, error) {
	conn, err := grpc.Dial(client.Address, grpc.WithInsecure())
	var result map[string]interface{}{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	c := pb.NewJwtServiceClient(conn)
	// 设置请求上下文的过期时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rs, err := c.VerifyJSON(ctx, &pb.VerifyRequest{Tokendata: token})
	if err != nil {
		return result, err
	}
	if rs.Status.Status == pb.StatusData_ERROR {
		return result, errors.New(rs.Status.Msg)
	}
	err = json.Unmarshal([]byte(rs.Payload), &result)
	return result, err
}

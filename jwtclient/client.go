package jwtclient

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	pb "github.com/Basic-Components/jwtrpc/jwtrpcdeclare"

	grpc "google.golang.org/grpc"
)

// Client jwt的客户端类型
type Client struct {
	Address string
}

// New 创建客户端对象
func New(address string) *Client {
	return &Client{Address: address}
}

// GetToken 获取签名
func (client *Client) GetToken(claims map[string]interface{}) (string, error) {
	conn, err := grpc.Dial(client.Address, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewJwtServiceClient(conn)
	// 设置请求上下文的过期时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	jsonBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	rs, err := c.Sign(ctx, &pb.SignRequest{Jsondata: jsonBytes})
	if err != nil {
		return "", err
	}
	if rs.Status.Status == pb.StatusData_ERROR {
		return "", errors.New(rs.Status.Msg)
	}
	return rs.Data, nil
}

// VerifyToken 验证签名
func (client *Client) VerifyToken(token string) (map[string]interface{}, error) {
	conn, err := grpc.Dial(client.Address, grpc.WithInsecure())
	var result map[string]interface{}
	if err != nil {
		return result, err
	}
	defer conn.Close()
	c := pb.NewJwtServiceClient(conn)
	// 设置请求上下文的过期时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rs, err := c.Verify(ctx, &pb.VerifyRequest{Tokendata: []byte(token)})
	if err != nil {
		return result, err
	}
	if rs.Status.Status == pb.StatusData_ERROR {
		return result, errors.New(rs.Status.Msg)
	}
	err = json.Unmarshal([]byte(rs.Data), &result)
	return result, err
}

package main

import (
	context "context"
	"encoding/json"
	config "github.com/Basic-Components/jwtrpc/config"
	pb "github.com/Basic-Components/jwtrpc/jwtrpcdeclare"
	"github.com/Basic-Components/jwtrpc/jwtsigner"
	"github.com/Basic-Components/jwtrpc/jwtverifier"
	"github.com/Basic-Components/jwtrpc/logger"
	"github.com/Basic-Components/jwtrpc/signals"
	"github.com/Basic-Components/jwtrpc/keygen"
	"net"

	grpc "google.golang.org/grpc"
	logrus "github.com/sirupsen/logrus"
)

type rpcserver struct {
	signer   *jwtsigner.Signer
	verifier *jwtverifier.Verifier
}

var signlog *logrus.Entry= logger.Log.WithFields(logrus.Fields{
	"API": "Signer",
})
var verifylog *logrus.Entry= logger.Log.WithFields(logrus.Fields{
	"API": "Verify",
})

func (s *rpcserver) Sign(ctx context.Context, in *pb.SignRequest) (*pb.SignResponse, error) {
	signlog.WithFields(logrus.Fields{
		"Jsondata": string(in.Jsondata),
	}).Debug("get request")
	result, err := s.signer.Sign(in.Jsondata)
	var status pb.StatusData
	if err != nil {
		status = pb.StatusData{
			Status: pb.StatusData_ERROR,
			Msg:    err.Error()}
		signlog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("send response")
		return &pb.SignResponse{
			Status: &status}, nil
	}
	status = pb.StatusData{
		Status: pb.StatusData_SUCCEED,
		Msg:    "success"}
	signlog.WithFields(logrus.Fields{
		"response": string(result),
	}).Debug("send response")
	return &pb.SignResponse{
		Status: &status,
		Data:   result}, nil
	
}
func (s *rpcserver) Verify(ctx context.Context, in *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	verifylog.WithFields(logrus.Fields{
		"tokendata": string(in.Tokendata),
	}).Debug("get request")
	result, err := s.verifier.Verify(in.Tokendata)
	var status pb.StatusData
	if err != nil {
		status = pb.StatusData{
			Status: pb.StatusData_ERROR,
			Msg:    err.Error()}
		verifylog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Debug("send response")
		return &pb.VerifyResponse{
			Status: &status}, nil
	} 
	status = pb.StatusData{
		Status: pb.StatusData_SUCCEED,
		Msg:    "success"}
	data, err := json.Marshal(result)
	if err != nil {
		status = pb.StatusData{
			Status: pb.StatusData_ERROR,
			Msg:    err.Error()}
		verifylog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Debug("send response")
		return &pb.VerifyResponse{
			Status: &status}, nil
	}
	verifylog.WithFields(logrus.Fields{
		"response": string(data),
	}).Debug("send response")
	return &pb.VerifyResponse{
		Status: &status,
		Data:   data}, nil
	
}

// Run 执行签名验签服务
func Run(conf config.ConfigType) {
	signer,err := jwtsigner.NewFromPath(
		conf.SignMethod,
		conf.PrivateKeyPath,
		map[string]interface{}{"iss": conf.Iss})
	if err != nil{
		logger.Log.Fatalf("init signer err: %v", err)
	}
	verifier, err := jwtverifier.NewFromPath(
		conf.SignMethod,
		conf.PublicKeyPath)
	if err != nil{
		logger.Log.Fatalf("init verifier err: %v", err)
	}
	rpc := rpcserver{signer: signer,verifier: verifier}
	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		logger.Log.Fatalf("failed to listen: %v", err)
		return
	}
	logger.Log.Info("server started @", conf.Address)
	server := grpc.NewServer()
	pb.RegisterJwtServiceServer(server, &rpc)
	if err := server.Serve(listener); err != nil {
		logger.Log.Fatalf("failed to serve: %v", err)
		return
	}
}

func main() {
	conf, signal := config.Init()
	if signal != nil {
		if signal == signals.GenkeySignal {
			keygen.AutoGenRsaKey()
		} else {
			logger.Log.Error("init config error %v", signal)
			return
		}
	} else {
		logger.Log.Info("start server config %v", conf)
		Run(conf)
	}
}
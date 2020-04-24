package main //import "github.com/Basic-Components/jwt-tools/jwtcenter"

import (
	context "context"
	"encoding/json"

	"net"

	pb "github.com/Basic-Components/jwt-tools/jwtcenter/jwtrpcdeclare"
	log "github.com/Basic-Components/jwt-tools/jwtcenter/logger"
	script "github.com/Basic-Components/jwt-tools/jwtcenter/script"
	"github.com/Basic-Components/jwt-tools/jwtsigner"
	"github.com/Basic-Components/jwt-tools/jwtverifier"
	"github.com/Basic-Components/jwtrpc/logger"
	grpc "google.golang.org/grpc"
)

type rpcserver struct {
	signer   *jwtsigner.Signer
	verifier *jwtverifier.Verifier
}

func (s *rpcserver) Sign(ctx context.Context, in *pb.SignRequest) (*pb.SignResponse, error) {
	log.Debug(map[string]interface{}{
		"in": string(in.),
		"Method":   "Sign",
	}, "get request")
	result, err := s.signer.Sign(in.Jsondata)
	var status pb.StatusData
	if err != nil {
		status = pb.StatusData{
			Status: pb.StatusData_ERROR,
			Msg:    err.Error()}
		log.Debug(map[string]interface{}{"error": err.Error(), "Method": "Sign"}, "send response")
		return &pb.SignResponse{
			Status: &status}, nil
	}
	status = pb.StatusData{
		Status: pb.StatusData_SUCCEED,
		Msg:    "success"}
	log.Debug(map[string]interface{}{"response": string(result), "Method": "Sign"}, "send response")
	return &pb.SignResponse{
		Status: &status,
		Data:   result}, nil

}
func (s *rpcserver) Verify(ctx context.Context, in *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	log.Debug(map[string]interface{}{"tokendata": string(in.Tokendata), "Method": "Verify"}, "get request")
	result, err := s.verifier.Verify(in.Tokendata)
	var status pb.StatusData
	if err != nil {
		status = pb.StatusData{
			Status: pb.StatusData_ERROR,
			Msg:    err.Error()}
		log.Debug(map[string]interface{}{"error": err.Error(), "Method": "Verify"}, "send response")
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
		log.Debug(map[string]interface{}{"error": err.Error(), "Method": "Verify"}, "send response")
		return &pb.VerifyResponse{
			Status: &status}, nil
	}
	log.Debug(map[string]interface{}{"response": string(data), "Method": "Verify"}, "send response")
	return &pb.VerifyResponse{
		Status: &status,
		Data:   data}, nil

}

// Run 执行签名验签服务
func Run(conf script.ConfigType) {
	signer, err := jwtsigner.NewFromPath(
		conf.SignMethod,
		conf.PrivateKeyPath,
		map[string]interface{}{"iss": conf.Iss})
	if err != nil {
		logger.Log.Fatalf("init signer err: %v", err)
	}
	verifier, err := jwtverifier.NewFromPath(
		conf.SignMethod,
		conf.PublicKeyPath)
	if err != nil {
		logger.Log.Fatalf("init verifier err: %v", err)
	}
	rpc := rpcserver{signer: signer, verifier: verifier}
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
	conf := script.Init()
	conf.ComponentName
	log.Init(conf.LogLevel, map[string]interface{}{"component_name": conf.ComponentName})
	log.Info(map[string]interface{}{"conf": conf}, "start server config %v")
	Run(conf)

}

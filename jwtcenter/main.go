package main //import "github.com/Basic-Components/jwt-tools/jwtcenter"

import (
	context "context"
	"os"

	"net"

	pb "github.com/Basic-Components/jwt-tools/jwtcenter/jwtrpcdeclare"
	log "github.com/Basic-Components/jwt-tools/jwtcenter/logger"
	script "github.com/Basic-Components/jwt-tools/jwtcenter/script"
	"github.com/Basic-Components/jwt-tools/jwtsigner"
	"github.com/Basic-Components/jwt-tools/jwtverifier"
	grpc "google.golang.org/grpc"
)

type rpcservice struct {
	AsymmetricSigner   *jwtsigner.Asymmetric
	AsymmetricVerifier *jwtverifier.Asymmetric
	SymmetricSigner    *jwtsigner.Symmetric
	SymmetricVerifier  *jwtverifier.Symmetric
}

// NewService 创建一个新的服务
func NewService(conf script.ConfigType) (*rpcservice, error) {
	s := &rpcservice{}
	asymmetricSigner, err := jwtsigner.AsymmetricFromPEMFile("RS256", conf.PrivateKeyPath)
	if err != nil {
		return s, err
	}
	asymmetricVerifier, err := jwtverifier.AsymmetricFromPEMFile("RS256", conf.PublicKeyPath)
	if err != nil {
		return s, err
	}
	symmetricSigner, err := jwtsigner.SymmetricNew("HS256", conf.Hashkey)
	if err != nil {
		return s, err
	}
	symmetricVerifier := jwtverifier.SymmetricNew(conf.Hashkey)
	if err != nil {
		return s, err
	}
	s.AsymmetricSigner = asymmetricSigner
	s.AsymmetricVerifier = asymmetricVerifier
	s.SymmetricSigner = symmetricSigner
	s.SymmetricVerifier = symmetricVerifier
	return s, nil
}

func (s *rpcservice) SignJSON(ctx context.Context, in *pb.SignJSONRequest) (*pb.SignResponse, error) {
	log.Debug(map[string]interface{}{
		"in":     in,
		"Method": "SignJSON",
	}, "get request")
	// result, err := s.signer.Sign(in.Jsondata)
	// var status pb.StatusData
	// if err != nil {
	// 	status = pb.StatusData{
	// 		Status: pb.StatusData_ERROR,
	// 		Msg:    err.Error()}
	// 	log.Debug(map[string]interface{}{"error": err.Error(), "Method": "Sign"}, "send response")
	// 	return &pb.SignResponse{
	// 		Status: &status}, nil
	// }
	// status = pb.StatusData{
	// 	Status: pb.StatusData_SUCCEED,
	// 	Msg:    "success"}
	// log.Debug(map[string]interface{}{"response": string(result), "Method": "Sign"}, "send response")
	return &pb.SignResponse{
		Status: &pb.StatusData{Status: pb.StatusData_SUCCEED},
		Token:  "123"}, nil

}
func (s *rpcservice) VerifyJSON(ctx context.Context, in *pb.VerifyRequest) (*pb.VerifyJSONResponse, error) {
	log.Debug(map[string]interface{}{"in": in, "Method": "VerifyJSON"}, "get request")
	// result, err := s.verifier.Verify(in.Tokendata)
	// var status pb.StatusData
	// if err != nil {
	// 	status = pb.StatusData{
	// 		Status: pb.StatusData_ERROR,
	// 		Msg:    err.Error()}
	// 	log.Debug(map[string]interface{}{"error": err.Error(), "Method": "Verify"}, "send response")
	// 	return &pb.VerifyResponse{
	// 		Status: &status}, nil
	// }
	// status = pb.StatusData{
	// 	Status: pb.StatusData_SUCCEED,
	// 	Msg:    "success"}
	// data, err := json.Marshal(result)
	// if err != nil {
	// 	status = pb.StatusData{
	// 		Status: pb.StatusData_ERROR,
	// 		Msg:    err.Error()}
	// 	log.Debug(map[string]interface{}{"error": err.Error(), "Method": "Verify"}, "send response")
	// 	return &pb.VerifyResponse{
	// 		Status: &status}, nil
	// }
	// log.Debug(map[string]interface{}{"response": string(data), "Method": "Verify"}, "send response")
	return &pb.VerifyJSONResponse{
		Status:  &pb.StatusData{Status: pb.StatusData_SUCCEED},
		Payload: `{"a":1}`}, nil
}

// Run 执行签名验签服务
func (s *rpcservice) Run() {
	listener, err := net.Listen("tcp", script.Config.Address)
	if err != nil {
		log.Logger.Fatalf("failed to listen: %v", err)
		return
	}
	log.Info(map[string]interface{}{"address": script.Config.Address}, "server started")
	server := grpc.NewServer()
	pb.RegisterJwtServiceServer(server, s)
	if err := server.Serve(listener); err != nil {
		log.Logger.Fatalf("failed to serve: %v", err)
		return
	}
}

func main() {
	err := script.Init()
	if err != nil {
		log.Warn(map[string]interface{}{"error": err}, "init config error")
		os.Exit(1)
	}
	log.Init(script.Config.LogLevel, map[string]interface{}{"component_name": script.Config.ComponentName})
	log.Info(map[string]interface{}{"conf": script.Config}, "setted server config")
	service, err := NewService(script.Config)
	if err != nil {
		log.Error(map[string]interface{}{"error": err}, "service not inited")
		os.Exit(1)
	}
	service.Run()
}

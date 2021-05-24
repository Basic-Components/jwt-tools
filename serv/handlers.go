package serv

import (
	"context"

	declare "github.com/Basic-Components/jwttools/jwtrpcdeclare"
	"github.com/Basic-Components/jwttools/options"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (s *Server) Meta(ctx context.Context, in *declare.MetaRequest) (*declare.MetaResponse, error) {
	return s.signer.Meta(), nil
}

func (s *Server) Sign(ctx context.Context, in *declare.SignRequest) (*declare.SignResponse, error) {
	payload := in.Payload
	opts := []options.SignOption{}
	if in.Sub != "" {
		opts = append(opts, options.WithSub(in.Sub))
	}
	if in.Aud != nil {
		for _, aud := range in.Aud {
			opts = append(opts, options.WithAud(aud))
		}
	}
	if in.Exp > 0 {
		opts = append(opts, options.WithExp(in.Exp))
	}
	if in.Nbf > 0 {
		opts = append(opts, options.WithNbf(in.Nbf))
	}

	res, err := s.signer.Sign(payload, opts...)
	if err != nil {
		return nil, err
	}
	return &declare.SignResponse{
		Token: res,
	}, nil
}

func (s *Server) Verify(ctx context.Context, in *declare.VerifyRequest) (*declare.VerifyResponse, error) {
	res, err := s.verifier.Verify(in.Token)
	if err != nil {
		return nil, err
	}
	resb, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return &declare.VerifyResponse{
		Payload: resb,
	}, nil
}

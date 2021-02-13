package endpoints

import (
	"authservice/service"
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthRequest)
		result, err := s.Login(ctx, req.Email, req.Password)
		response := AuthResponse{Token: result, Success: 1, Message: "Login Successful!"}
		if err != nil {
			response.Success = 0
			response.Message = err.Error()
		}
		return response, nil
	}
}

func makeSignUpEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(NewUserRequest)
		err := s.SignUp(ctx, req.NewUser)
		response := AuthResponse{Token: "", Success: 1, Message: "SignUp Successful!"}
		if err != nil {
			response.Success = 0
			response.Message = err.Error()
		}
		return response, nil
	}
}

func makeLoginWithGoogleEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context,request interface{}) (interface{},error) {
		req := request.(OAuth2Request)
		result, err := s.LoginWithGoogle(ctx,req.Token)
		response := AuthResponse{Token: result, Success: 1, Message: "Login Successful!"}
		if err != nil {
			response.Success = 0
			response.Message = err.Error()
		}
		return response, nil
	}
}
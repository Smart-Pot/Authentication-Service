package endpoints

import (
	"authservice/data"
	"authservice/service"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Login  endpoint.Endpoint
	SignUp endpoint.Endpoint
}

type AuthResponse struct {
	Token   string
	Success int32
	Message string
}

type AuthRequest struct {
	Email    string
	Password string
}

type NewUserRequest struct {
	NewUser data.SignUpForm
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Login:  makeLoginEndpoint(s),
		SignUp: makeSignUpEndpoint(s),
	}
}

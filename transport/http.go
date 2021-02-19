package transport

import (
	"authservice/data"
	"authservice/endpoints"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const userIDTag = "x-user-id"

func MakeHTTPHandlers(e endpoints.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter().PathPrefix("/auth").Subrouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		e.Login,
		decodeAuthHTTPRequest,
		encodeHTTPResponse,
		options...,
	))
	
	r.Methods("POST").Path("/google/login").Handler(httptransport.NewServer(
		e.LoginWithGoogle,
		decodeOAuthHTTPRequest,
		encodeHTTPResponse,
		options...,
	))

	r.Methods("POST").Path("/signup").Handler(httptransport.NewServer(
		e.SignUp,
		decodeNewUserHTTPRequest,
		encodeHTTPResponse,
		options...,
	))

	r.Methods("POST").Path("/verify/{hash}").Handler(httptransport.NewServer(
		e.Verify,
		decodeVerifyRequest,
		encodeHTTPResponse,
		options...
	))

	r.Methods("GET").Path("/").Handler(httptransport.NewServer(
		e.Resolve,
		decodeAuthHTTPRequest,
		encodeResolveHTTPResponse,
		options...,
	))

	return r
}

func encodeHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeResolveHTTPResponse(ctx context.Context,w http.ResponseWriter,response interface{}) error {
	a := response.(endpoints.AuthResponse)
	w.Header().Set("X-Forwarded-User",a.Token)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return nil
}

func decodeOAuthHTTPRequest(_ context.Context,r *http.Request) (interface{}, error) {
	const tn = "x-oauth-token"
	t := r.Header.Get(tn) 
	return endpoints.OAuth2Request{
		Token: t,
	},nil
}


func decodeVerifyRequest(_ context.Context,r *http.Request) (interface{},error) {
	vars := mux.Vars(r)
	h := vars["hash"]
	return endpoints.VerifyRequest{
		Hash : h,
	},nil

}

func decodeResolveHTTPRequest(_ context.Context, r *http.Request) (interface{},error) {
	jwt := r.Header.Get("x-auth-token")
	if jwt == "" {
		return nil,errors.New("token not found")
	}
	return endpoints.OAuth2Request {
		Token: jwt,
	},nil
}

func decodeAuthHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var ep  struct{
		Email string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&ep); err != nil {
		return nil, err
	}

	return endpoints.AuthRequest{
		Email:    ep.Email,
		Password: ep.Password,
	}, nil
}

func decodeNewUserHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var u data.SignUpForm

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, err
	}
	return endpoints.NewUserRequest{
		NewUser: u,
	}, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	fmt.Println("FO:UN: ERROR  http.go:127")
	w.WriteHeader(403)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

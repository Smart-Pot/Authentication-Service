package transport

import (
	"authservice/data"
	"authservice/endpoints"
	"context"
	"encoding/json"
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

	r.Methods("POST").Path("/signup").Handler(httptransport.NewServer(
		e.SignUp,
		decodeNewUserHTTPRequest,
		encodeHTTPResponse,
		options...,
	))

	return r
}

func encodeHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeAuthHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var u data.UserCredentials

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, err
	}

	return endpoints.AuthRequest{
		Email:    u.Email,
		Password: u.Password,
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
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

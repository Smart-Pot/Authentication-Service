package transport

import (
	"authservice/data"
	"authservice/endpoints"
	"context"
	"encoding/json"
	"net/http"

	"github.com/Smart-Pot/pkg/common/constants"
	"github.com/Smart-Pot/pkg/common/perrors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	errTokenNotFound = perrors.New("token not found", 400)
)

func MakeHTTPHandlers(e endpoints.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter().PathPrefix("/auth").Subrouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(perrors.EncodeHTTPError),
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
		options...,
	))

	r.Methods("GET").Path("/").Handler(httptransport.NewServer(
		e.Resolve,
		decodeResolveHTTPRequest,
		encodeResolveHTTPResponse,
		options...,
	))

	return handlers.CORS(headers, methods, origins)(r)
}

func encodeHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeResolveHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	a := response.(endpoints.AuthResponse)
	w.Header().Set(constants.UserIDHeaderKey, a.Token)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return nil
}

func decodeOAuthHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	t := r.Header.Get(constants.OAuthHeaderKey)
	return endpoints.OAuth2Request{
		Token: t,
	}, nil
}

func decodeVerifyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	h := vars["hash"]
	return endpoints.VerifyRequest{
		Hash: h,
	}, nil

}

func decodeResolveHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	jwt := r.Header.Get(constants.TokenHeaderKey)
	if jwt == "" {
		return nil, perrors.New("token not found in header", http.StatusBadRequest)
	}
	return endpoints.OAuth2Request{
		Token: jwt,
	}, nil
}

func decodeAuthHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var ep struct {
		Email    string
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

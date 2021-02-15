// Package service implements a service for Authentication service for: Login, Signup, Forgot password operations.
package service

import (
	"authservice/data"
	"authservice/service/oauth"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Smart-Pot/jwtservice"
	"github.com/Smart-Pot/pkg/adapter/amqp"
	"github.com/Smart-Pot/pkg/tool/crypto"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrEmailTaken codes returned by failures when email is already is taken
	ErrEmailTaken = errors.New("Email is already taken")
	// ErrWrongPassword codes returned by failures for wrong password
	ErrWrongPassword = errors.New("Password is wrong for given email")
	// ErrEmailNotFound codes returned when there is no user with the given email
	ErrEmailNotFound = errors.New("Not found an account by given email")
	// ErrNotVerified codes returned when the user tries to log in inactive account
	ErrNotVerified = errors.New("Account not verified")
)

type service struct {
	logger   log.Logger
	jwt      *jwtservice.JwtService
	producer amqp.Producer
}
// Service represents an authentication service 
type Service interface {
	LoginWithGoogle(ctx context.Context,token string)(string,error)
	Login(ctx context.Context, email, password string) (string,error)
	SignUp(ctx context.Context, form data.SignUpForm) error
	Verify(ctx context.Context,hash string) error	
}

// NewService creates a new service for given parameters
func NewService(logger log.Logger, producer amqp.Producer) Service {
	return &service{
		logger:   logger,
		jwt : jwtservice.New(),
		producer: producer,
	}
}

// SignUp gets a form data and verify it and notify amqp server
func (s *service) SignUp(ctx context.Context, form data.SignUpForm) error {
	// Try to find a user who has that email
	_, err := data.GetUserCrediantals(ctx, form.Email)
	// If a cred is founded than return email taken error
	if err == nil {
		return ErrEmailTaken
	}

	// if err is not credental not found, return 'err'	
	if err != nil && err != data.ErrCredentalNotFound {
		return err
	}

	if err := form.HashPassword(); err != nil {
		return err
	}
	form.GenerateUserID()

	if err = data.CreateUser(ctx,form); err != nil {
		return err
	}

	// Hash user id for verification mail
	h, err := crypto.Encrypt(form.UserID)
	if err != nil {
		return err
	}

	r := struct {
		Hash  string `json:"hash"`
		Email string `json:"email"`
	}{
		Hash:  h,
		Email: form.Email,
	}

	b, _ := json.Marshal(r)

	return s.producer.Produce(b)	
}

// Login gets email and password, and generate JWT for userId
func (s *service) Login(ctx context.Context, email, password string) (string,error) {
	var token string
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "Login",
			"param:email", email,
			"param:password", password,
			"result", token,
			"took", time.Since(beginTime))
	}(time.Now())

	userCred, err := data.GetUserCrediantals(ctx, email)
	if err != nil {
		if err == data.ErrCredentalNotFound {
			return "" , ErrEmailNotFound
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userCred.Password), []byte(password)); err != nil {
		return "", ErrWrongPassword
	}
	token,err = s.jwt.Tokenize(userCred.UserID)

	return token,err
}


func (s *service) LoginWithGoogle(ctx context.Context,token string) (string,error) {
	var result string
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "Login",
			"param:token", token,
			"result", result,
			"took", time.Since(beginTime))
	}(time.Now())

	claim,err := oauth.ValidateGoogleJWT(token)
	if err != nil {
		return "",err
	}
	cred, err := data.GetUserCrediantals(ctx, claim.Email)
	
	// If user exist, tokenize userId and returns it
	if err == nil {		
		return s.jwt.Tokenize(cred.UserID)
	}

	f := data.SignUpForm{
		Email: claim.Email,
		FirstName: claim.FirstName,
		LastName: claim.LastName,
		Password: "",
		IsOAuth: true,
	}
	if err = f.HashPassword(); err != nil {
		return "", err
	}
	f.GenerateUserID()
	

	
	if err = data.CreateUser(ctx,f);  err != nil {
		return "",err
	}

	result,err = s.jwt.Tokenize(f.UserID)
	return result,err
}


func (s service) Verify(ctx context.Context, hash string) error {
	var err error
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "Login",
			"param:token", hash,
			"result", err,
			"took", time.Since(beginTime))
	}(time.Now())
	id, err := crypto.Decrypt(hash)
	if err != nil {
		return err
	}
	if err := data.UpdateUserRecord(ctx, id, "active", true); err != nil {
		return err
	}
	return nil
}

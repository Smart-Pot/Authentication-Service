// Package service implements a service for Authentication service for: Login, Signup, Forgot password operations.
package service

import (
	"authservice/data"
	"authservice/service/oauth"
	"context"
	"encoding/json"
	"errors"

	"github.com/Smart-Pot/jwtservice"
	"github.com/Smart-Pot/pkg/adapter/amqp"
	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrEmailTaken codes returned by failures when email is already is taken
	ErrEmailTaken = errors.New("Email is already taken")
	// ErrWrongPassword codes returned by failures for wrong password
	ErrWrongPassword = errors.New("Password is wrong for given email")
	// ErrEmailNotFound codes returned when there is no user with the given email
	ErrEmailNotFound = errors.New("Not found an account by given email")
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
}

// NewService creates a new service for given parameters
func NewService(logger log.Logger, producer amqp.Producer) Service {
	return &service{
		logger:   logger,
		producer: producer,
		jwt : jwtservice.New(),
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

	return s.signUp(form)
}

// Login gets email and password, and generate JWT for userId
func (s *service) Login(ctx context.Context, email, password string) (string,error) {
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
	return s.jwt.Tokenize(userCred.UserID)
}


func (s *service) LoginWithGoogle(ctx context.Context,token string) (string,error) {
	
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
	f.GenerateUserID()

	if err := s.signUp(f); err != nil {
		return "",err
	}

	return s.jwt.Tokenize(f.UserID)
}


func (s *service) signUp(form data.SignUpForm) error {
	// Stringfy form data
	b, err := json.Marshal(form)
	if err != nil {
		return err
	}

	// Notify rabbitmq server
	if err := s.producer.Produce(b); err != nil {
		return err
	}
	return nil
}
 

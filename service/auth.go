package service

import (
	"authservice/data"
	"context"
	"errors"

	"github.com/Smart-Pot/jwtservice"
	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger log.Logger
	jwt    *jwtservice.JwtService
}

type Service interface {
	Login(ctx context.Context, email, password string) (token string, err error)
	SignUp(ctx context.Context, form data.SignUpForm) error
}

func NewService(logger log.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (s *service) SignUp(ctx context.Context, form data.SignUpForm) error {
	_, err := data.GetUserCrediantals(ctx, form.Email)
	if err != data.ErrCredentalNotFound {
		if err == nil {
			return errors.New("Email is already taken")
		}
		return err
	}
	if err := form.HashPassword(); err != nil {
		return err
	}

	form.GenerateUserID()

	// TODO: notify message broker server for new user

	return nil
}

func (s *service) Login(ctx context.Context, email, password string) (token string, err error) {
	userCred, err := data.GetUserCrediantals(ctx, email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userCred.Password), []byte(password)); err != nil {
		return "", errors.New("wrong password")
	}
	return s.jwt.Tokenize(userCred.UserId)
}

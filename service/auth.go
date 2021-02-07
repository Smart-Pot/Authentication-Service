package service

import (
	"authservice/data"
	"context"
	"encoding/json"
	"errors"

	"github.com/Smart-Pot/jwtservice"
	"github.com/Smart-Pot/pkg/adapter/amqp"
	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger   log.Logger
	jwt      *jwtservice.JwtService
	producer amqp.Producer
}

type Service interface {
	Login(ctx context.Context, email, password string) (token string, err error)
	SignUp(ctx context.Context, form data.SignUpForm) error
}

func NewService(logger log.Logger, producer amqp.Producer) Service {
	return &service{
		logger:   logger,
		producer: producer,
	}
}

func (s *service) SignUp(ctx context.Context, form data.SignUpForm) error {
	// Verify email is not taken
	_, err := data.GetUserCrediantals(ctx, form.Email)

	if err == nil {
		return errors.New("Email is already taken")
	}

	if err != nil && err != data.ErrCredentalNotFound {
		return err
	}

	if err := form.HashPassword(); err != nil {
		return err
	}

	form.GenerateUserID()

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

func (s *service) Login(ctx context.Context, email, password string) (token string, err error) {
	userCred, err := data.GetUserCrediantals(ctx, email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userCred.Password), []byte(password)); err != nil {
		return "", errors.New("wrong password")
	}
	return s.jwt.Tokenize(userCred.UserID)
}

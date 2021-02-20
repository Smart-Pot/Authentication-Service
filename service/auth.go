// Package service implements a service for Authentication service for: Login, Signup, Forgot password operations.
package service

import (
	"authservice/data"
	"authservice/service/oauth"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Smart-Pot/pkg/adapter/amqp"
	"github.com/Smart-Pot/pkg/common/perrors"
	"github.com/Smart-Pot/pkg/tool/crypto"
	"github.com/Smart-Pot/pkg/tool/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrEmailTaken codes returned by failures when email is already is taken
	ErrEmailTaken = perrors.New("Email is already taken",http.StatusBadRequest)
	// ErrWrongPassword codes returned by failures for wrong password
	ErrWrongPassword = perrors.New("Password is wrong for given email",http.StatusForbidden)
	// ErrEmailNotFound codes returned when there is no user with the given email
	ErrEmailNotFound = perrors.New("Not found an account by given email",http.StatusForbidden)
	// ErrInactiveAccount codes returned when the user tries to log in inactive account
	ErrInactiveAccount = perrors.New("Email address not verified",http.StatusForbidden)
	// ErrInvalidPwd :
	ErrInvalidPwd = perrors.New("Password is not valid",http.StatusBadRequest)
	// ErrInvalidHash :
	ErrInvalidHash = perrors.New("Hash can not be decrypted",http.StatusBadRequest)
	// ErrInvalidToken :
	ErrInvalidToken = perrors.New("Token can not be resolved", http.StatusBadRequest)
	// ErrInvalidEmail : 
	ErrInvalidEmail = perrors.New("Email is not valid",http.StatusBadRequest)
	// 
	errServer = perrors.FromStatusCode(http.StatusInternalServerError)
)

type service struct {
	logger   log.Logger
	verifyProducer amqp.Producer
	forgotProducer amqp.Producer
	verificationTimeout time.Duration
}
// Service represents an authentication service 
type Service interface {
	LoginWithGoogle(ctx context.Context,token string)(string,error)
	Login(ctx context.Context, email, password string) (string,error)
	SignUp(ctx context.Context, form data.SignUpForm) error
	Verify(ctx context.Context,hash string) error	
	Resolve(ctx context.Context,token string) (*jwt.AuthToken,error)
}

// NewService creates a new service for given parameters
func NewService(logger log.Logger, verifyProducer,forgotProducer amqp.Producer) Service {
	return &service{
		logger:   logger,
		forgotProducer: forgotProducer,
		verifyProducer : verifyProducer,
		verificationTimeout: 48*time.Hour,
	}
}

// SignUp gets a form data and verify it and notify amqp server
func (s *service) SignUp(ctx context.Context, form data.SignUpForm) error {
	u, err := data.GetUserByEmail(ctx, form.Email)
    // If a cred is founded than return email taken error
    if err == nil {
		const layout = "2021-02-08 01:02:15.0271274 +0000 UTC"
        ct, _ := time.Parse(layout, u.Date)
        if u.Active || (time.Since(ct) < s.verificationTimeout) {
            return ErrEmailTaken
        }
        data.RemoveUser(ctx, u.ID)
    }


	// if err is not credental not found, return 'err'	
	if err != nil && err != data.ErrUserNotFound {
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
	h, err := crypto.VerifyMailCip.Encrypt(form.UserID)
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

	return s.verifyProducer.Produce(b)	
}

func (s *service) Resolve(ctx context.Context,tokenStr string) (*jwt.AuthToken,error) {
	var err error
	var at *jwt.AuthToken
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "Resolve",
			"param:email", tokenStr,
			"result:err",err,
			"result:tkn",at,

			"took", time.Since(beginTime))
	}(time.Now())

	at,err =jwt.Verify(tokenStr)
	
	if err == jwt.ErrTokenExpired {
		return nil,perrors.FromError("",400,err)
	}
	if err != nil {
		return nil,ErrInvalidToken
	}
	return at,nil
}	

func (s *service) ForgotPassword(ctx context.Context,email string)  error {
	var err error
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "ForgorPassword",
			"param:email", email,
			"result:err",err,
			"took", time.Since(beginTime))
	}(time.Now())
	_,err = data.GetUserByEmail(ctx,email)
	if err != nil {
		return errServer
	}

	if !data.ValidateEmail(email) {
		return ErrInvalidEmail
	}

	h,err := crypto.ForgotPwdCip.Encrypt(email)
	if err != nil {
		return errServer
	}

	if err = s.forgotProducer.Produce([]byte(h)); err != nil {
		return errServer
	}
	return nil
}

func (s *service) UpdatePassword(ctx context.Context,hash,newPwd string) error {
	var err error
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "UpdatePassword",
			"param:hash", hash,
			"result:err",err,
			"took", time.Since(beginTime))
	}(time.Now())
	
	email, err := crypto.ForgotPwdCip.Decrypt(hash);
	if  err != nil {
		if err == data.ErrUserNotFound {
			return ErrEmailNotFound
		}
		return errServer
	}

	u,err := data.GetUserByEmail(ctx,email)

	if err != nil {
		return errServer
	}

	if !data.ValidatePassword(newPwd) {
		return ErrInvalidPwd
	}

	f := data.SignUpForm{}
	f.Password = newPwd
	if err = f.HashPassword(); err != nil {
		return perrors.FromStatusCode(http.StatusInternalServerError)
	}

	if err = data.UpdateUserRecord(ctx,u.ID,"password",f.Password); err  != nil {
		return errServer
	}

	return nil
}

// Login gets email and password, and generate JWT for userId
func (s *service) Login(ctx context.Context, email, password string) (string,error) {
	var token string
	var err error
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "Login",
			"param:email", email,
			"param:password", password,
			"result:token", token,
			"result:err",err,
			"took", time.Since(beginTime))
	}(time.Now())

	u, err := data.GetUserByEmail(ctx, email)
	if err != nil {
		if err == data.ErrUserNotFound {
			return "" , ErrEmailNotFound
		}
		return "", errServer
	}

	if !u.Active {
		return "", ErrInactiveAccount
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", ErrWrongPassword
	}

	tkn := &jwt.AuthToken {
		UserID: u.ID,
		DeviceIDs: u.Devices,
		Authorization: u.Authorization,
	}

	token,err = jwt.Tokenize(tkn)

	return token,errServer
}


func (s *service) LoginWithGoogle(ctx context.Context,token string) (string,error) {
	var result string
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "LoginWithGoogle",
			"param:token", token,
			"result", result,
			"took", time.Since(beginTime))
	}(time.Now())

	claim,err := oauth.ValidateGoogleJWT(token)
	if err != nil {
		return "", errServer
	}
	u, err := data.GetUserByEmail(ctx, claim.Email)
	
	// If user exist, tokenize userId and returns it
	if err == nil {	
		tkn := &jwt.AuthToken{
			UserID: u.ID,
			DeviceIDs: u.Devices,
			Authorization: u.Authorization,
		}	
		return jwt.Tokenize(tkn)
	}

	// if error is not a notfound error then returns it
	if err != nil && err != data.ErrUserNotFound {
		return "",errServer
	}

	f := data.SignUpForm{
		Email: claim.Email,
		FirstName: claim.FirstName,
		LastName: claim.LastName,
		Password: "",
		IsOAuth: true,
	}
	if err = f.HashPassword(); err != nil {
		return "", ErrInvalidPwd
	}
	f.GenerateUserID()
	
	if err = data.CreateUser(ctx,f);  err != nil {
		return "",errServer
	}
	tkn := &jwt.AuthToken{
		UserID: f.UserID,
		DeviceIDs: []string{""},
		Authorization: 0,
	}	
	result,err = jwt.Tokenize(tkn)
	if err != nil {
		return "", ErrInvalidToken
	}
	return result,nil
}


func (s service) Verify(ctx context.Context, hash string) error {
	var err error
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "Verify",
			"param:token", hash,
			"result", err,
			"took", time.Since(beginTime))
	}(time.Now())
	id, err := crypto.VerifyMailCip.Decrypt(hash)
	if err != nil {
		return ErrInvalidHash
	}
	if err := data.UpdateUserRecord(ctx, id, "active", true); err != nil {
		if err == data.ErrUserNotFound {
			return ErrInvalidHash
		}
		return errServer
	}
	return nil
}

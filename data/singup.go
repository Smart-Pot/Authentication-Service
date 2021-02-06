package data

import (
	"regexp"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpForm struct {
	UserID    string
	FirstName string `validate:"is-name,required"`
	LastName  string `validate:"is-name,required"`
	Email     string `validate:"email,required"`
	Password  string `validate:"is-passwd,required"`
}

func (s *SignUpForm) Validate() error {
	v := validator.New()
	v.RegisterValidation("is-name", validateName)
	v.RegisterValidation("is-passwd", validateName)

	return v.Struct(s)
}

func (s *SignUpForm) HashPassword() error {
	newPasswd, err := hashPassword(s.Password)
	if err != nil {
		return err
	}
	s.Password = newPasswd
	return nil
}

func (s *SignUpForm) GenerateUserID() {
	s.UserID = uuid.NewString()
}

var nameValidationReg = regexp.MustCompile("^[a-zA-Z,ç,Ç,ğ,Ğ,ı,İ,ö,Ö,ş,Ş,ü,Ü]*$")
var passwordValidationReg = regexp.MustCompile("^[a-zA-Z,ç,Ç,ğ,Ğ,ı,İ,ö,Ö,ş,Ş,ü,Ü]*$")

func validateName(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return nameValidationReg.MatchString(s)
}

func validatePassword(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return passwordValidationReg.MatchString(s)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

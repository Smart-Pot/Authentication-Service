package data

import (
	"regexp"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type SignUpForm struct {
	FirstName string `json:"firstname" validate:"is-name,required"`
	LastName  string `json:"lastname"  validate:"is-name,required"`
	Email     string `json:"email"     validate:"email,required"`
	Password  string `json:"password"  validate:"is-passwd,required"`
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

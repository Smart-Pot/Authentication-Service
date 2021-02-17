package data

import (
	"regexp"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	symbols = `\.\*\_`
)

var (
	nameRegexp = regexp.MustCompile("^[a-zA-Z,ç,Ç,ğ,Ğ,ı,İ,ö,Ö,ş,Ş,ü,Ü]*$")
	upperCaseRegexp = regexp.MustCompile("[A-Z]")
	lowerCaseRegexp = regexp.MustCompile("[a-z]")
	digitRegexp = regexp.MustCompile("[0-9]")
	symbolRegexp = regexp.MustCompile("["+symbols+"]") // TODO add more symbols 
	passwdRegexp = regexp.MustCompile("^["+symbols+"A-Za-z0-9]*.{6,}$")
)

type SignUpForm struct {
	UserID    string
	FirstName string `json:"firstname" validate:"is-name,required"`
	LastName  string `json:"lastname"  validate:"is-name,required"`
	Email     string `json:"email"     validate:"email,required"`
	Password  string `json:"password"  validate:"is-passwd,required"`
	IsOAuth   bool `json:"isOAuth" bson:"isOAuth"`
}

func (s *SignUpForm) Validate() error {
	v := validator.New()
	v.RegisterValidation("is-name", nameValidator)
	v.RegisterValidation("is-passwd", passwordValidator)

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


func nameValidator(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return nameRegexp.MatchString(s)
}

func passwordValidator(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return ValidatePassword(s)
}



func ValidatePassword(pwd string)bool {
	if !passwdRegexp.MatchString(pwd) {
		return false
	}
	if !upperCaseRegexp.MatchString(pwd) {
		return false
	}
	if !lowerCaseRegexp.MatchString(pwd) {
		return false
	}
	if !digitRegexp.MatchString(pwd) {
		return false
	}
	if !symbolRegexp.MatchString(pwd) {
		return false
	}
	return true
}
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

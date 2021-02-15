package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUpValidation(t *testing.T) {
	cases := []struct {
		form SignUpForm
		ok bool
	}{
		{
			form: SignUpForm{
				FirstName: "Ahmet",
				LastName: "ÖZCAN",
				Password: "ahmetcanozcan",
				Email: "ahmetcanozcan7@gmail.com",
			},
			ok: true,
		},
		{
			form: SignUpForm{
				FirstName: "K€N@N",
				LastName: "@BB@C",
				Password: "kenanabbak",
				Email: "kenanabbak@gmail.com",
			},
			ok: true,},
		}

	for _, c := range cases {
		err := c.form.Validate()
		if c.ok{
			assert.Nil(t,err)
		} else {
			assert.NotNil(t,err)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct{
		s string
		ok bool
	}{
		{"Sifre123",false},
		{"Sifre.123",true},
		{"sifre123",false},
		{"Sifre__123",true},
		{"sifre.",false},
		{"Ss.1",false},
	}
	for _,ti := range tests {
		assert.Equal(t,ti.ok,validatePassword(ti.s),"For "+ti.s)
	}

}
func TestS(t *testing.T) {
	assert.Equal(t,true,symbolRegexp.MatchString("Sifre.123"))

}
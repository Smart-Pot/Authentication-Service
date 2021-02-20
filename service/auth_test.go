package service

import (
	"authservice/data"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Smart-Pot/pkg"
	"github.com/Smart-Pot/pkg/db"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

const inactive = "inactive_"

var (
	_s Service
	_todo = context.TODO()
	_user   = data.User{
		ID: "test_user_id",
		Email: "testuser@testmail.com",
		Password: "testuserpassword",
		Active: true,
	}
)



// mockProducer reperesents pkg/adapter/amqp.Producer for testing
type mockProducer struct{}

func (p *mockProducer) Produce(b []byte) error {
	return nil
}
 
func TestMain(m *testing.M) {

	l := log.NewJSONLogger(ioutil.Discard)
	_s = NewService(l,&mockProducer{},&mockProducer{})

	wd,_ := os.Getwd()
	pkg.ConfigOptions.BaseDir = filepath.Join(wd,"..","config")
	if err := pkg.Config.ReadConfig(); err != nil {
		panic(err)
	}

	if err := db.Connect(db.PkgConfig("users")); err != nil {
		panic(err)
	}

	p := _user.Password // Store original password
	if err := _user.HashPassword(); err !=nil{
		panic(err)
	}
	if err := data.SaveUser(_todo,_user); err != nil {
		panic(err)
	}
	// Inactive User
	inactiveUser := data.User{
		ID: inactive+_user.ID,
		Email: inactive+_user.Email,
		Password: _user.Password,
		Active: false,
	}
	if err := data.SaveUser(_todo,inactiveUser); err != nil {
		panic(err)
	}
	_user.Password = p
	inactiveUser.Password = p
	c := m.Run()

	if err := data.RemoveUser(_todo,_user.ID); err != nil {
		panic(err)
	}
	if err := data.RemoveUser(_todo,inactiveUser.ID); err != nil {
		panic(err)
	}

	os.Exit(c)
}

func TestService_Login(t *testing.T) {
	tests := []struct{
		email string
		password string
		err error
	}{
		{
			email:_user.Email,
			password: _user.Password,
			err: nil,
		},
		{
			email:inactive+_user.Email,
			password: _user.Password,
			err: ErrInactiveAccount,
		},
		{
			email: _user.Email,
			password: "wrong_password",
			err : ErrWrongPassword,
		},
		{
			email: "wrongmail@wmail.com",
			password: "",
			err : ErrEmailNotFound,
		},
	}

	for _,test := range tests {
		token,err := _s.Login(_todo,test.email,test.password)
		assert.Equal(t,test.err,err)
		if test.err != nil {
			assert.Equal(t,"",token)
		} else {
			assert.NotEqual(t,"",token) // Token is not empty
			/*
				Note:
					Token can be validated by decoding using github.com/Smart-Pot/jwtservice
					but it exceeds purpose of the test.	
			*/
		}
	}
	
}



func TestService_SignUp(t *testing.T) {
	tests := []struct{
			err error
			form data.SignUpForm
		
		}{
			{
				err: ErrEmailTaken,
				form : data.SignUpForm{
				UserID: "test_user_id",
				Email: "testuser@testmail.com",
				Password: "testuserpassword",
				FirstName: "test",
				LastName: "user",
			  	},
			},
			{
				err: nil,
				form : data.SignUpForm{
				UserID: "test_user_id",
				Email: "nottaken@testmail.com",
				Password: "testuserpassword",
				FirstName: "test",
				LastName: "user",
			  	},
			},
		}

	for _,test := range tests {
		err := _s.SignUp(_todo,test.form)
		assert.Equal(t,err,test.err)
	}	
}




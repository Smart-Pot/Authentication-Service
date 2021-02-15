package data

import (
	"context"
	"errors"
	"time"

	"github.com/Smart-Pot/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


var (
	ErrUserNotFound = errors.New("User not found")
)

type User struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	Password      string   `json:"password"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	Image         string   `json:"image"`
	Date          string   `json:"date"`
	Active        bool     `json:"active"`
	Authorization int      `json:"authorization"`
	Devices       []string `json:"devices"`
	OAuth 		  bool     `json:"oauth"`	
}


func (u *User) HashPassword() error {
	pwd, err := hashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}


func CreateUser(ctx context.Context,form SignUpForm) error {
	u := User {
		ID : form.UserID,
		FirstName: form.FirstName,
		LastName: form.LastName,
		Email: form.LastName,
		Password: form.Password,
		Date: time.Now().UTC().String(),
		Image: "",
		Devices: nil,
		Authorization: 0,
		Active: false,
		OAuth: form.IsOAuth,
	}
	_, err := db.Collection().InsertOne(ctx, u)
	return err
}

func GetUserByEmail(ctx context.Context,email string) (*User,error) {
	r := db.Collection().FindOne(ctx, bson.M{"email": email})
	var u User
	if err := r.Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil,ErrUserNotFound
		}
		return nil,err
	}
	return &u,nil
}

func UpdateUserRecord(ctx context.Context, id, key string, value interface{}) error {
	filter := bson.M{"id": id}

	updateUser := bson.M{"$set": bson.M{
		key: value,
	}}

	res, err := db.Collection().UpdateOne(ctx, filter, updateUser)

	if res.ModifiedCount <= 0 {
		return errors.New("record can not updated")
	}

	if err != nil {
		return err
	}

	return nil
}



/* These two functions are using for testing. */


func SaveUser(ctx context.Context,user User) error {
	_, err := db.Collection().InsertOne(ctx,user)
	return err
}

func RemoveUser(ctx context.Context,userId string) error {
	res,err :=db.Collection().DeleteOne(ctx,bson.M{"id":userId})
	if res.DeletedCount == 0 {
		return errors.New("no doc deleted")
	}
	return err
}


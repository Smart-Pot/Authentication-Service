package data

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCredentials struct {
	UserID   string `json:"userId" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (cred *UserCredentials) HashPassword() error {
	newPasswd, err := hashPassword(cred.Password)
	if err != nil {
		return err
	}
	cred.Password = newPasswd
	return nil
}


func GetUserCrediantals(ctx context.Context, email string) (*UserCredentials, error) {
	res := collection.FindOne(ctx, bson.M{"email": email})
	var cred UserCredentials
	err := res.Decode(&cred)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCredentalNotFound
		}
		return nil, err
	}
	return &cred, nil
}


/* These two functions are using for testing. */


func SaveUserCrediantals(ctx context.Context,cred UserCredentials) error {
	_, err := collection.InsertOne(ctx,cred)
	return err
}

func RemoveUserCrediantals(ctx context.Context,userId string) error {
	res,err :=collection.DeleteOne(ctx,bson.M{"id":userId})
	if res.DeletedCount == 0 {
		return errors.New("no doc deleted")
	}
	return err
}


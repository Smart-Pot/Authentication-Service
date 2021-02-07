package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCredentials struct {
	UserId           string `json:"userId"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	IsActive         string `json:"isActive"`
	Authorization    string `json:"authorization"`
	VerificationHash string `json:"verificationHash"`
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

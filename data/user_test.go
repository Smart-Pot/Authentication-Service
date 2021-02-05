package data

import (
	"authservice/config"
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestNofFoundError(t *testing.T) {
	config.ReadConfig()
	DatabaseConnection()
	_, err := GetUserCrediantals(context.TODO(), "asdasd")
	if err == mongo.ErrNoDocuments {
		return
	}
	if err != nil {
		t.Error("ERR!", err)
		t.FailNow()
	}
}

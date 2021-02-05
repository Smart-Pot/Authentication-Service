package service

import "testing"

func TestResolve(t *testing.T) {
	n := "Ahmet"
	s := &jwtService{}
	tokenStr, err := s.tokenize(n)
	if err != nil {
		t.Error("Err", err)
		t.FailNow()
	}
	claims, err := s.verify(tokenStr)
	if err != nil {
		t.Error("Err", err)
		t.FailNow()
	}
	if claims["userId"] != n {
		t.Error("Mismatch claim.UserID", claims["userId"], "n", n)
		t.FailNow()
	}
}

func TestGetUserID(t *testing.T) {
	n := "Ahmet"
	s := &jwtService{}
	tokenStr, err := s.tokenize(n)
	if err != nil {
		t.Error("Err", err)
		t.FailNow()
	}
	id, err := s.getUserID(tokenStr)
	if err != nil {
		t.Error("Err", err)
		t.FailNow()
	}
	if id != n {
		t.Error("ERR! id", id, "n", n)
		t.FailNow()
	}
}

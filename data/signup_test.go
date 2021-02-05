package data

import (
	"testing"
)

func TestSignUpValidation(t *testing.T) {
	cases := []struct {
		FirstName string
		res       bool
	}{
		{
			FirstName: "Ahmet",
			res:       true,
		},
		{
			FirstName: "Ken@n",
			res:       false,
		},
		{
			FirstName: "123Ahmet1",
			res:       false,
		},
		{
			FirstName: "Ege",
			res:       true,
		},
		{
			FirstName: "Ege CEMAL",
			res:       false,
		},
	}

	for i, c := range cases {
		s := SignUpForm{
			FirstName: c.FirstName,
		}
		err := s.Validate()
		t.Log("err->", err, "res", c.res, "case", i)
		if (err == nil) != c.res {
			t.Error("Failed on case", i, "validation err", err, "is ok", c.res)
			t.FailNow()
		}
	}
}

package jwt_test

import (
	"golang/advanced/pkg/jwt"
	"testing"
)

func TestJWT_Create(t *testing.T) {
	const email = "aa2@aaa.com"
	jwtService := jwt.NewJWT("3LqePlmvOhs7hZs6t1W2YLRMz8qFJc2uP17JIXD6IGA=")
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is InValid")
	}
	if data.Email != email {
		t.Fatalf("Email %s is not equal to %s", data.Email, email)
	}
}

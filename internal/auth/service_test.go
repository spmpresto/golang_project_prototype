package auth_test

import (
	"golang/advanced/internal/auth"
	"golang/advanced/internal/user"
	"testing"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "aa2@aaa.com",
	}, nil
}

func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "aa2@aaa.com"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err := authService.Register(initialEmail, "1", "Vasya")
	if err != nil {
		t.Fatal(err)
	}
	if email != initialEmail {
		t.Fatalf("Email %s do not match %s", email, initialEmail)
	}
}

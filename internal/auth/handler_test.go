package auth_test

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"golang/advanced/configs"
	"golang/advanced/internal/auth"
	"golang/advanced/internal/user"
	"golang/advanced/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}

	return &handler, mock, nil
}

func TestLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("aa2@aaa.com", "$2a$10$ik7SdFh/FbYoE50lwT/9kup3vA97nCQEL5Fikn0WI/rujkO/r3Z.K")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	if err != nil {
		t.Fatal(err)
		return
	}
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "aa2@aaa.com",
		Password: "1",
	})
	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/login", reader)
	if err != nil {
		t.Fatal(err)
	}
	handler.Login().ServeHTTP(wr, req)
	if wr.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, wr.Result().StatusCode)
	}

}

func TestHandlerRegisterSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}

	// FindByEmail should return no rows (user does not exist)
	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	// Expect Create (INSERT)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "new@test.com",
		Password: "password123",
		Name:     "Test User",
	})
	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/register", reader)
	if err != nil {
		t.Fatal(err)
	}
	handler.Register().ServeHTTP(wr, req)
	if wr.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, wr.Result().StatusCode)
	}
}

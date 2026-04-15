package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"golang/advanced/internal/auth"
	"golang/advanced/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "aa2@aaa.com",
		Password: "$2a$10$ik7SdFh/FbYoE50lwT/9kup3vA97nCQEL5Fikn0WI/rujkO/r3Z.K",
		Name:     "Vasya",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "aa2@aaa.com").
		Delete(&user.User{})

}

func TestLoginSuccess(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "aa2@aaa.com",
		Password: "1",
	})
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Response code is %v", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resData auth.LoginResponse
	if err := json.Unmarshal(body, &resData); err != nil {
		t.Fatal(err)
	}
	if resData.Token == "" {
		t.Fatalf("Token is empty")
	}

	removeData(db)
}

func TestLoginFail(t *testing.T) {
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "aa2@aaa.com",
		Password: "2",
	})
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 401 {
		t.Fatalf("Response code is %v", res.StatusCode)
	}
	removeData(db)
}

package main

import (
	"github.com/joho/godotenv"
	"golang/advanced/internal/link"
	"golang/advanced/internal/stat"
	"golang/advanced/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
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
	db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
}

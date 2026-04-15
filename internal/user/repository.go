package user

import (
	"golang/advanced/pkg/db"
)

type UserRepository struct {
	database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.database.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.database.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

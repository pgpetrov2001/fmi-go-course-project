package dao

import (
	"course-project/entities"
	"course-project/utils"
	"errors"
	"gorm.io/gorm"
)

type CPNS struct {
	Db *gorm.DB
}

func (cpns *CPNS) CreateUser(email string, username string, pass string) error {
	hash, err := utils.HashPassword(pass)
	if err != nil {
		panic(err)
	}
	err = cpns.Db.Create(&entities.User{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
	}).Error
	return err
}

func (cpns *CPNS) Authenticate(email string, pass string) (bool, error) {
	var user entities.User
	err := cpns.Db.First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return utils.CheckPasswordHash(pass, user.PasswordHash), nil
}

func (cpns *CPNS) GetUsers() ([]entities.User, error) {
	var users []entities.User
	err := cpns.Db.Find(&users).Error
	return users, err
}

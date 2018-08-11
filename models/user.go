package models

import (
	"regexp"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model

	Username  string  `gorm:"not null",json:"username"`
	Password  string  `json:"password"`
}

func GetUser(username string) (user User, err error) {
	err = db.Where("username=?", username).First(&user).Error

	return
}

func SetUser(username string, password string) (user User, err error) {
	// Hash password if not already a bcrypt hash
	reBCrypt, err := regexp.Compile(`^\$2[ayb]\$.{56}$`)
	if err != nil {return}

	if !reBCrypt.MatchString(password) {
		password, _ = HashPassword(password)
	}

	newUser := User{Username: username, Password: password}

	logger.Debugf("New User: %s", newUser.Username)

	err = db.Create(&newUser).Error

	if err != nil {
		logger.Errorf("Error creating user %s: %s", newUser.Username, err)
	}

	return
}


func UserCount() int64 {
	var count int64
	db.Model(&User{}).Count(&count)

	return count
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
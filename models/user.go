package models

import (
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint       `json:"id" gorm:"primary_key"`

	Username  string     `json:"username" gorm:"not null"`
	Password  string     `json:"-"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// GetID client id
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func GetUser(id int) (user User, err error) {
	err = db.Where("id=?", id).First(&user).Error

	return
}

func GetUserByUsername(username string) (user User, err error) {
	err = db.Where("username=?", username).First(&user).Error

	return
}

func SetUser(usr User) (user User, err error) {
	// Hash password if not already a bcrypt hash
	reBCrypt, err := regexp.Compile(`^\$2[ayb]\$.{56}$`)
	if err != nil {return}

	if !reBCrypt.MatchString(usr.Password) {
		usr.Password, _ = hashPassword(usr.Password)
	}

	err = db.Create(&usr).Error
	if err != nil {
		logger.Errorf("Error creating user %s: %s", usr.Username, err)
	}

	return
}

func UserCount() int64 {
	var count int64
	db.Model(&User{}).Count(&count)

	return count
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

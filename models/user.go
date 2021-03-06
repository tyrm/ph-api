package models

import (
	"regexp"
	"time"

	"github.com/google/jsonapi"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint       `jsonapi:"attr,db_id" gorm:"primary_key"`

	Username  string     `jsonapi:"primary,user" gorm:"not null"`
	Password  string     `jsonapi:"attr,password,omitempty" gorm:"-"`
	PasswordHash  string `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func GetUser(id int) (user User, err error) {
	err = db.Where("id=?", id).First(&user).Error

	return
}

func GetUserByUsername(username string) (user User, err error) {
	err = db.Where("lower(username) = lower(?)", username).First(&user).Error

	return
}

func GetUsersPage(count int, page int) (users []User, err error) {
	offset := count * page;
	err = db.Limit(count).Offset(offset).Find(&users).Error

	return
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *User) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}

func SetUser(usr User) (user User, err error) {
	// Hash password if not already a bcrypt hash
	reBCrypt, err := regexp.Compile(`^\$2[ayb]\$.{56}$`)
	if err != nil {return}

	if reBCrypt.MatchString(usr.Password) {
		usr.PasswordHash = usr.Password
	} else {
		usr.PasswordHash, _ = hashPassword(usr.Password)
	}

	err = db.Create(&usr).Error
	if err != nil {
		logger.Errorf("Error creating user %s: %s", usr.Username, err)
	}

	user = usr

	return
}

func UserCount() int64 {
	var count int64
	db.Model(&User{}).Count(&count)

	return count
}

func UserExists(username string) (exists bool, err error) {
	var count int64
	err = db.Model(&User{}).Where("lower(username) = lower(?)", username).Count(&count).Error

	if count > 0 {
		exists = true
	} else {
		exists = false
	}

	return
}

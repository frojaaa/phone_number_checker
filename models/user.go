package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"os"
	"phone_numbers_checker/errors"
	"strings"
)

type User struct {
	gorm.Model
	Username               string `gorm:"size:255;not null;unique;" json:"username"`
	Password               string `gorm:"size:255;not null;" json:"password"`
	EnteredCheckerPassword bool   `gorm:"not null;" json:"enteredCheckerPassword"`
}

func (u User) SaveUser() (User, error) {
	fmt.Println("Saving user ", u.Username)
	u.BeforeSave()
	fmt.Println(u.Password)
	err := DB.Create(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (u User) UpdateUser(field string, value any) (User, error) {
	fmt.Println("Saving user ", u.Username)
	u.BeforeSave()
	fmt.Println(u.Password)
	err := DB.Model(User{}).Where("username = ?", u.Username).Update(field, value).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	errors.HandleError("Error while hashing password", &err)

	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.EnteredCheckerPassword = false
}

func VerifyPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (User, error) {
	u := User{}

	err := DB.Model(User{}).Where("username = ?", username).Take(&u).Error
	if err != nil {
		return User{}, err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return User{}, err
	}

	return u, nil
}

func CheckerPasswordCheck(username string, checkerPassword string) (User, error) {
	u, err := GetUserByUsername(username)
	err = VerifyPassword(checkerPassword, os.Getenv("CHECKER_PASSWORD_HASH"))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return User{}, err
	}
	return u, nil
}

func GetUserByUsername(username string) (User, error) {
	u := User{}

	err := DB.Model(User{}).Where("username = ?", username).Take(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

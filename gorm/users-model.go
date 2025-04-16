package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func CreateUser(db *gorm.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	fmt.Println("User created successfully")
	return nil
}

func Login(db *gorm.DB, user *User) (string, error) {
	selectedUser := new(User)
	result := db.Where("email = ?", user.Email).First(selectedUser)
	if result.Error != nil {
		return "", result.Error
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(selectedUser.Password),
		[]byte(user.Password),
	)

	if err != nil {
		return "", err
	}

	jwtSecretKey := "your-secret-key"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = selectedUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

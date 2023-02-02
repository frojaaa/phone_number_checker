package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"phone_numbers_checker/errors"
	"phone_numbers_checker/models"
)

type UserAuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CheckerPasswordInput struct {
	Password string `json:"password" binding:"required"`
}

var userKey = "user"

func Login(c *gin.Context) {
	session := sessions.Default(c)
	sessionUser := session.Get(userKey)
	fmt.Println(sessionUser)
	if sessionUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logout first"})
		return
	}
	var input UserAuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	u.Username = input.Username
	u.Password = input.Password

	token, err := models.LoginCheck(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password is incorrect."})
		return
	}
	session.Set(userKey, u.Username)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error at server"})
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logout first"})
		return
	}
	var input UserAuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	u, err := u.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session.Set(userKey, u.Username)
	err = session.Save()
	errors.HandleError("Error while saving session: ", &err)
	c.JSON(http.StatusCreated, gin.H{"msg": "Вы успешно зарегистрировались"})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	log.Println("logging out user:", user)
	if user == nil {
		log.Println("Invalid session token")
		return
	}
	session.Delete(userKey)
	if err := session.Save(); err != nil {
		log.Println("Failed to save session: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Logged out successfully"})

}

func CheckerPassword(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(userKey).(string)
	if username == "" {
		log.Println("Invalid session token")
		return
	}
	var input CheckerPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	password := input.Password
	user, err := models.CheckerPasswordCheck(username, password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Wrong password"})
		return
	}
	user, err = user.UpdateUser("entered_checker_password", true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Entered checker pwd successfully!"})
}

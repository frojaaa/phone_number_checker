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

var userKey = "user"

func Login(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	fmt.Println(user)
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

	err := models.LoginCheck(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username or password is incorrect."})
		return
	}
	session.Set(userKey, u.Username)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Error at server"})
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Authorized successfully"})
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

	_, err := u.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session.Set(userKey, u.Username)
	err = session.Save()
	errors.HandleError("Error while saving session: ", &err)
	c.JSON(http.StatusCreated, gin.H{"msg": "Вы успешно зарегистрировались", "username": u.Username})
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
		log.Println("Failed to save session:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Logged out successfully"})

}

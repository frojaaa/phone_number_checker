package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

func Login(c *gin.Context) {
	var input UserAuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	u.Username = input.Username
	u.Password = input.Password

	u, token := models.LoginCheck(u.Username, u.Password)
	fmt.Println("Entered checker pwd: ", u.EnteredCheckerPassword)
	c.JSON(http.StatusOK, gin.H{"token": token, "enteredCheckerPassword": u.EnteredCheckerPassword})
}

func Register(c *gin.Context) {
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

	errors.HandleError("Error while saving session: ", &err)
	c.JSON(http.StatusCreated, gin.H{"msg": "Вы успешно зарегистрировались"})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Logged out successfully"})

}

func CheckerPassword(c *gin.Context) {
	var input UserAuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	password := input.Password
	user, err := models.CheckerPasswordCheck(input.Username, password)
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

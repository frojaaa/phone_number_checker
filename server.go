package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"phone_numbers_checker/controllers"
)

func RunServer() http.Handler {
	router := gin.Default()
	router.POST("/run", controllers.RunChecker)
	return router
}

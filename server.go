package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"phone_numbers_checker/controllers"
	"phone_numbers_checker/errors"
	"phone_numbers_checker/middleware"
	"phone_numbers_checker/models"
)

func RunServer() http.Handler {
	models.ConnectDB()

	router := gin.Default()

	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(os.Getenv("STORE_SECRET")))
	errors.HandleError("Error while creating redis store: ", &err)
	router.Use(sessions.Sessions("session", store))

	protected := router.Group("/checker")
	protected.Use(middleware.AuthRequired)
	protected.POST("/run", controllers.RunChecker)
	protected.POST("/checkPassword", controllers.CheckerPassword)

	public := router.Group("/")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	public.POST("/logout", controllers.Logout)
	return router
}

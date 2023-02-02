package models

import (
	"fmt"
	"phone_numbers_checker/errors"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	env, err := godotenv.Read()
	errors.HandleError("Error while reading .env: ", &err)

	host := env["DB_HOST"]
	user := env["DB_USER"]
	password := env["DB_PASSWORD"]
	port := env["DB_PORT"]
	name := env["DB_NAME"]

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, name)
	fmt.Println(DBURL)

	DB, err = gorm.Open(mysql.Open(DBURL))
	errors.HandleError("Error while connecting to DB: ", &err)
	fmt.Println("Connected to DB")
	err = DB.AutoMigrate(&User{})
	errors.HandleError("Error while migration: ", &err)
}

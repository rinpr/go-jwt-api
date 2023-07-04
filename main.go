package main

import (
	"fmt"
	AuthController "rinpr/jwt-gin-api/controller/auth"
	UserController "rinpr/jwt-gin-api/controller/user"
	"rinpr/jwt-gin-api/middleware"
	"rinpr/jwt-gin-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	orm.InitDB()
	initENV()

	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	authorized := r.Group("/users", middleware.JWTAuthen())
	authorized.GET("/readall", UserController.ReadAll)
	authorized.GET("/profile", UserController.Profile)

  	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func initENV() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file!")
	}
}
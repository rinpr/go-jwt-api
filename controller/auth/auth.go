package auth

import (
	"fmt"
	"net/http"
	"os"
	"rinpr/jwt-gin-api/orm"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var hmacSampleSecret []byte

func Register(c *gin.Context) {
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user in database
	var userExists orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExists)
	if userExists.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"message": "This username is already exists!",
		})
		return
	}

	// Create user and save to database
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := orm.User{Username: json.Username, Password: string(encryptedPassword), Fullname: json.Fullname, Avatar: json.Avatar}
	orm.Db.Create(&user)
	if (user.ID > 0) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "User registered successfully!",
			"userId": user.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"message": "User registered failed!",
		})
	}
}

func Login(c *gin.Context) {

	// Check if returned json has both Username and Password
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user in database
	var userExists orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExists)
	if userExists.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"message": "This username doesn't exists!",
		})
		return
	}

	// Validate password
	err := bcrypt.CompareHashAndPassword([]byte(userExists.Password), []byte(json.Password))
	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExists.ID,
			"exp": time.Now().Add(time.Minute * 1).Unix(),
		})
		// Sign and get the complete encoded token as a string using secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Login successfully!",
			"token": tokenString,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"message": "Invalid password!",
		})
	}
}
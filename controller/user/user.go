package user

import (
	"net/http"
	"rinpr/jwt-gin-api/orm"

	"github.com/gin-gonic/gin"
)

var hmacSampleSecret []byte

func ReadAll(c *gin.Context) {
	var users []orm.User
	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "User read successfully!",
		"users": users,
	})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user orm.User
	orm.Db.First(&user, userId)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "User read successfully!",
		"users": user,
	})
}
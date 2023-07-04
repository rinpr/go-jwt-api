package orm

import (
	"fmt"
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

var Db *gorm.DB

func InitDB() {
	dsn := "root@tcp(127.0.0.1:3306)/gojwt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error: Failed to connect to database!")
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	Db = db
}
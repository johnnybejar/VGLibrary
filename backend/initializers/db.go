package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDBConn() {
	username := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	address := os.Getenv("DBADDRESS")
	port := os.Getenv("DBPORT")
	dbname := os.Getenv("DBNAME")

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, port, dbname)
	fmt.Println(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to db")
	}
}
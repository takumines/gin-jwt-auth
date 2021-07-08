package db

import (
	"github.com/takumines/gin-jwt-auth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := "host=127.0.0.1 user=test password=test1234 dbname=go_auth port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	connect, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = connect

	connect.AutoMigrate(&models.User{}, &models.PasswordReset{})
	if err != nil {
		panic(err)
	}
}

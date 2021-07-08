package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takumines/gin-jwt-auth/db"
	"github.com/takumines/gin-jwt-auth/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	data := map[string]string{}

	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if data["password"] != data["password_confirm"] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Passwords do not match!",
		})
		return
	}

	// パスワードをエンコードする
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	db.DB.Create(&user)

	c.JSON(http.StatusOK, user)
}

func Login(c *gin.Context) {
	data := map[string]string{}

	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user := models.User{}

	db.DB.Where("email = ?", data["email"]).First(&user)

	// emailが存在しない場合
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	// passwordが一致しない場合
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Incorrect password",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}
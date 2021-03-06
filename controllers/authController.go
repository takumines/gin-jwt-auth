package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/takumines/gin-jwt-auth/db"
	"github.com/takumines/gin-jwt-auth/models"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.StandardClaims
}

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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect password",
		})
		return
	}

	// JWT
	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Server Error")
		return
	}

	// Cookie
	c.SetCookie("jwt", token, time.Now().Add(time.Hour*24).Second(), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"jwt": token,
	})
}

func User(c *gin.Context) {
	// CookieからJWTを取得
	cookie, _ := c.Cookie("jwt")
	// tokenを取得
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
	}
	claims := token.Claims.(*Claims)
	// User IDを取得
	user := models.User{}
	db.DB.Where("id = ?", claims.Issuer).First(&user)

	c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt", "", time.Now().Add(-time.Hour).Second(), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

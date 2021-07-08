package controllers

import (
	"math/rand"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"
	"github.com/takumines/gin-jwt-auth/db"
	"github.com/takumines/gin-jwt-auth/models"
	"golang.org/x/crypto/bcrypt"
)

func Forget(c *gin.Context) {
	data := map[string]string{}

	// リクエストデータを格納する
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}
	db.DB.Create(&passwordReset)

	from := "example@test.com"
	to := []string{
		data["email"],
	}
	url := "http://localhost:3000/reset/" + token
	message := []byte("Click <a href=\"" + url + "\">here</a> to reset your password")
	err := smtp.SendMail("localhost:1025", nil, from, to, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func Reset(c *gin.Context) {
	data := map[string]string{}

	// リクエストデータを格納する
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// パスワードチェック
	if data["password"] != data["password_confirm"] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password do not match!",
		})
		return
	}

	passwordReset := models.PasswordReset{}

	// tokenからデータを取得し、一番直近のデータを取得してtokenチェックを行う
	if err := db.DB.Where("token = ?", data["token"]).Last(&passwordReset); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token!",
		})
		return
	}

	//パスワードをエンコード
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	db.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	// パスワードリセットに使用したtokenと一致するpasswortResetモデルを全て削除する
	db.DB.Exec("delete from password_resets where token = ?", passwordReset.Token)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// ランダムな文字列を返す関数
func RandStringRunes(n int) string {
	lettersRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = lettersRunes[rand.Intn(len(lettersRunes))]
	}
	return string(b)
}

package routes

import (
	"net/http"
	"rumahsakit/db"
	"rumahsakit/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
			return
		}

		user.Password = string(hashedPassword)

		if err := db.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal register user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Registrasi berhasil", "user": user.Username})
	})

	r.POST("/login", func(c *gin.Context) {
		var input models.User
		var user models.User

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Cek user berdasarkan username
		if err := db.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
			return
		}

		// Cek password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Login berhasil",
			"user_id":  user.ID,
			"username": user.Username,
		})
	})
}

package handlers

import (
	"encoding/json"
	"loginapp/database"
	"loginapp/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("supersecretkey")

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Register (otomatis role = "user")
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Input tidak valid", http.StatusBadRequest)
		return
	}

	// Cek apakah username sudah ada
	var exists int64
	database.DB.Model(&models.User{}).Where("username = ?", user.Username).Count(&exists)
	if exists > 0 {
		http.Error(w, "Username sudah terdaftar", http.StatusBadRequest)
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal memproses password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashed)
	user.Role = "user"

	// Simpan ke DB
	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Gagal mendaftar", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Registrasi berhasil!"})
}

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Input tidak valid", http.StatusBadRequest)
		return
	}

	var user models.User
	err := database.DB.Where("username = ?", input.Username).First(&user).Error
	if err != nil {
		http.Error(w, "Username tidak ditemukan", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Password salah", http.StatusUnauthorized)
		return
	}

	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Gagal membuat token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"role":  user.Role,
	})
}

// Protected (contoh endpoint yang butuh token & role admin)
func Protected(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token tidak ditemukan", http.StatusUnauthorized)
		return
	}

	// Biasanya header format: "Bearer tokenstring", jadi kita split
	tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	if tokenStr == "" {
		http.Error(w, "Token tidak valid", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token tidak valid", http.StatusUnauthorized)
		return
	}

	if claims.Role != "admin" {
		http.Error(w, "Akses hanya untuk admin", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Akses berhasil untuk admin!",
		"role":    claims.Role,
	})
}

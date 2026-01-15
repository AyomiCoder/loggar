package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

var db *sql.DB

// SetDB sets the database connection for handlers
func SetDB(database *sql.DB) {
	db = database
}

// LoginHandler handles user authentication
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Query user from database
	var userID int
	var passwordHash string
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", req.Email).
		Scan(&userID, &passwordHash)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(userID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Save token to database
	_, err = db.Exec("INSERT INTO tokens (user_id, token) VALUES ($1, $2)", userID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		Email: req.Email,
	})
}

// generateJWT creates a JWT token for the user
func generateJWT(userID int, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-me"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24 * 365).Unix(), // 1 year expiry
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

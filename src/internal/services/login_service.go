package services

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go_api/internal/models"
	"encoding/hex"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type LoginService struct {
	db     *gorm.DB
	jwtKey []byte
}

func NewLoginService(db *gorm.DB) *LoginService {
	// Resolve path to root relative to current file
	rootPath, err := filepath.Abs(filepath.Join("./"))
	if err != nil {
		log.Fatal("Unable to resolve root path:", err)
	}

	envPath := filepath.Join(rootPath, ".env")

	log.Println("Loaded .env from:", envPath)
	// Load .env from root
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatal("JWT_KEY is not set in the environment variables")
	}
	log.Printf("==> JWT_KEY value: %s", jwtKey)

	return &LoginService{
		db:     db,
		jwtKey: []byte(jwtKey),
	}
}

func (s *LoginService) Authenticate(username, password string) (string, error) {
	var user models.User

	log.Printf("Authenticate: Attempting to authenticate user: %s", username)

	// Fetch the user from the database by username, preloading the Role
	if err := s.db.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Authenticate: Error fetching user '%s' from DB: %v", username, err)
		return "", errors.New("invalid credentials")
	}

	if user.ID == 0 { // Check if user was actually found.  Important!
		log.Printf("Authenticate: User '%s' not found in database", username)
		return "", errors.New("invalid credentials")
	}

	log.Printf("Authenticate: Fetched user from DB: Username: %s", user.Username)
    log.Printf("Authenticate: Fetched user from DB: Hashed Password (string): %s", user.Password)
    log.Printf("Authenticate: Fetched user from DB: Hashed Password (bytes): %s", hex.EncodeToString([]byte(user.Password)))


	// Compare the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Authenticate: Password mismatch for user '%s': %v", username, err)
		return "", errors.New("invalid credentials")
	}

	log.Printf("Authenticate: Password comparison successful for user: %s", username)

	// Create JWT claims
	claims := &Claims{
		Username: username,
		Role:     user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		log.Printf("Authenticate: Error generating JWT token for user '%s': %v", username, err)
		return "", errors.New("failed to generate token")
	}

	log.Printf("Authenticate: Token generated successfully for user: %s", username)

	// Log the successful login attempt
	login := models.Login{
		Username:  username,
		Password:  "", // Do not store the password in the login log
		LoginTime: time.Now(),
	}
	if err := s.db.Create(&login).Error; err != nil {
		log.Println("Authenticate: Error logging login attempt:", err)
		// Do not return this error to the client as the login was successful
	}

	log.Printf("Authenticate: Successful login for user: %s", username)
	return tokenString, nil
}



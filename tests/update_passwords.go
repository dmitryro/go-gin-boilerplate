package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ocr_api/internal/services" // Adjust the import path if necessary
)

func main() {
	// Load environment variables (adjust path if needed)
	rootPath, err := filepath.Abs(filepath.Join("./")) // Assuming this script is in the same directory as main.go
	if err != nil {
		log.Fatal("Unable to resolve root path:", err)
	}
	envPath := filepath.Join(rootPath, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Database connection setup
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable client_encoding=UTF8", //Added client_encoding
		os.Getenv("PG_HOST"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DATABASE"), os.Getenv("PG_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userService := services.NewUserService(db) // Create an instance of your UserService

	usersToUpdate := map[string]string{
		"admin":    "nu45edi1",
	}

	for username, password := range usersToUpdate {
		err := userService.UpdateUserPassword(username, password)
		if err != nil {
			log.Printf("Error updating password for user '%s': %v", username, err)
		} else {
			log.Printf("Successfully updated password for user '%s'", username)
		}
	}

	fmt.Println("Password update script finished.")
}

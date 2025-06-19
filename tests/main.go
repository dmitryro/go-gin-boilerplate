package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	plaintext := "nu45edi1"
	hashed := "$2a$10$3/3A8R4z30HWGJj7bUAjkO/7BBwoVedibX3/jGYtSif4x83d6NwLa"

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	if err != nil {
		log.Printf("bcrypt.CompareHashAndPassword FAILED in bcrypt_test.go: %v %s", err, hashed)
		fmt.Println("bcrypt comparison failed")
	} else {
		log.Println("bcrypt.CompareHashAndPassword SUCCEEDED in bcrypt_test.go")
		fmt.Println("bcrypt comparison succeeded")
	}

	// Test with a known good hash
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt.GenerateFromPassword failed: %v", err)
	}
	hashedString := string(hashedBytes) // Convert []byte to string
	err = bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(plaintext))
	if err != nil {
		log.Printf("bcrypt.CompareHashAndPassword FAILED with generated hash: %v", err)
		fmt.Println("bcrypt comparison failed with generated hash")
	} else {
		log.Println("bcrypt.CompareHashAndPassword SUCCEEDED with generated hash %s", hashedString)
		fmt.Println("bcrypt comparison succeeded with generated hash")
	}
}


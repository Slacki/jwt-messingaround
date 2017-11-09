package main

import (
	"bio-cleaner-api/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	
)

func migrate() {
	db.DropTableIfExists(&models.User{})

	db.CreateTable(&models.User{})
}

func fakeData() {
	pass, err := bcrypt.GenerateFromPassword([]byte("test"), 14)
	if err != nil {
		log.Fatal("Bcrypt hashing unsuccessful!")
	}
	db.Create(&models.User{
		Identifier:   "test",
		PasswordHash: string(pass),
		Email:        "test@test.test",
	})
}

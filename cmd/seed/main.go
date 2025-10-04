package main

import (
	"log"
	"os"

	"github.com/findwannawhy/IAD-Semester5/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	sql, err := os.ReadFile("resources/db/seed.sql")
	if err != nil {
		log.Fatalf("failed to read seed file: %v", err)
	}

	if err := db.Exec(string(sql)).Error; err != nil {
		log.Fatalf("seeding failed: %v", err)
	}

	log.Println("seeding: OK")
}

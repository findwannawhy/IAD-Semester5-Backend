package main

import (
	"log"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
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

	// 1) Базовая схема
	if err := db.AutoMigrate(
		&ds.User{},           // users
		&ds.AcidSolubleSample{},       // samples (образцы)
		&ds.ImpurityFractionExperiment{},     // experiments (расчет массовой доли примесей)
		&ds.ExperimentSample{}, // experiments_samples (m-m эксперименты—образцы)
	); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}


	log.Println("migrations: OK")
}
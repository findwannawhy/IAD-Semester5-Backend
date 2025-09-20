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
		&ds.Material{},       // materials (услуги)
		&ds.Experiment{},     // experiments (заявки)
		&ds.ExperimentItem{}, // experiment_items (m-m заявки—услуги)
	); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	// 2) Доп. ограничения, которых не всегда делает AutoMigrate "как надо"
	//    — проверка допустимых статусов заявки,
	//    — частичный уникальный индекс "один черновик на пользователя",
	//    — на всякий случай: составной уникальный индекс для m-m (идемпотентно).

	sql := []string{
		// CHECK по статусам чтобы статус был только из списка
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'chk_experiments_status'
			) THEN
				ALTER TABLE experiments
				ADD CONSTRAINT chk_experiments_status
				CHECK (status IN ('draft','deleted','formed','finished','rejected'));
			END IF;
		END$$;`,
	
		// Один черновик на пользователя
		`CREATE UNIQUE INDEX IF NOT EXISTS uid_one_draft_per_user
			 ON experiments (creator_id)
		 WHERE status = 'draft';`,
	
		// Составной уникальный ключ для m-m
		`CREATE UNIQUE INDEX IF NOT EXISTS uid_experiment_material
			 ON experiment_items (experiment_id, material_id);`,
	}

	for _, s := range sql {
		if err := db.Exec(s).Error; err != nil {
			log.Fatalf("migration step failed: %v\nSQL: %s", err, s)
		}
	}

	log.Println("migrations: OK")
}
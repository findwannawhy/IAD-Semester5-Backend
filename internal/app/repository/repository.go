package repository

import (
	"errors"

	minioClient "github.com/findwannawhy/IAD-Semester5/internal/app/minio"
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotAllowed    = errors.New("not allowed")
	ErrNoDraft       = errors.New("no draft for this user")
)

type Repository struct {
	db *gorm.DB
	mc *minio.Client
	userId uint
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	mc, err := minioClient.InitMinio()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		mc: mc,
		userId: 0,
	}, nil
}

func (r *Repository) GetUserID() (uint) {
	return r.userId
}

func (r *Repository) SetUserID(id uint) {
	r.userId = id
}

func (r *Repository) SignOut() {
	r.userId = 0
}
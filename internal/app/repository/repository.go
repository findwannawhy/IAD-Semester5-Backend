package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	minioClient "github.com/findwannawhy/IAD-Semester5/internal/app/minio"

	"github.com/go-redis/redis"
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
	rd *redis.Client
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

	rd := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &Repository{
		db: db,
		mc: mc,
		rd: rd,
	}, nil
}


func (r *Repository) GetToken(userID string) (string, error) {
	token, err := r.rd.Get(userID).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func blacklistKeyForToken(tokenString string) string {
	h := sha256.Sum256([]byte(tokenString))
	return "blacklist:" + hex.EncodeToString(h[:])
}

func (r *Repository) AddTokenToBlacklist(ctx context.Context, tokenString string, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}
	key := blacklistKeyForToken(tokenString)
	return r.rd.Set(key, "1", ttl).Err()
}

func (r *Repository) IsTokenBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	key := blacklistKeyForToken(tokenString)
	n, err := r.rd.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Гостевые сессии

func guestSessionKey(sessionID string) string {
	return "guest_session:" + sessionID
}

func guestSessionViewedKey(sessionID string) string {
	return "guest_viewed:" + sessionID
}

// CreateGuestSession создает новую гостевую сессию в Redis
func (r *Repository) CreateGuestSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	key := guestSessionKey(sessionID)
	// Сохраняем время создания сессии
	return r.rd.Set(key, time.Now().Unix(), ttl).Err()
}

// GetGuestSessionCreatedAt возвращает время создания гостевой сессии
func (r *Repository) GetGuestSessionCreatedAt(ctx context.Context, sessionID string) (time.Time, error) {
	key := guestSessionKey(sessionID)
	timestamp, err := r.rd.Get(key).Int64()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(timestamp, 0), nil
}

// DeleteGuestSession удаляет гостевую сессию из Redis
func (r *Repository) DeleteGuestSession(ctx context.Context, sessionID string) error {
	sessionKey := guestSessionKey(sessionID)
	viewedKey := guestSessionViewedKey(sessionID)
	
	// Удаляем и сессию, и список просмотренных образцов
	pipe := r.rd.Pipeline()
	pipe.Del(sessionKey)
	pipe.Del(viewedKey)
	_, err := pipe.Exec()
	return err
}

// AddViewedSample добавляет просмотренный образец в список для гостевой сессии
func (r *Repository) AddViewedSample(ctx context.Context, sessionID string, sampleID uint, ttl time.Duration) error {
	key := guestSessionViewedKey(sessionID)
	
	// Выполняем все операции в одном pipeline
	pipe := r.rd.Pipeline()
	// Сначала удаляем все вхождения этого образца (если есть)
	pipe.LRem(key, 0, sampleID)
	// Добавляем ID образца в начало списка
	pipe.LPush(key, sampleID)
	// Оставляем последние 4 просмотренных образца (т.к. текущий будет исключен на фронте)
	pipe.LTrim(key, 0, 3)
	// Обновляем TTL
	pipe.Expire(key, ttl)
	
	_, err := pipe.Exec()
	return err
}

// GetViewedSamples возвращает список ID просмотренных образцов для гостевой сессии
func (r *Repository) GetViewedSamples(ctx context.Context, sessionID string) ([]uint, error) {
	key := guestSessionViewedKey(sessionID)
	
	values, err := r.rd.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	
	sampleIDs := make([]uint, 0, len(values))
	for _, v := range values {
		var id uint
		_, err := fmt.Sscanf(v, "%d", &id)
		if err == nil {
			sampleIDs = append(sampleIDs, id)
		}
	}
	
	return sampleIDs, nil
}
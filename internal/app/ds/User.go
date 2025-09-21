package ds

import (
	"time"
)

type User struct {
	ID          int       `gorm:"primaryKey"                   json:"id"`
	Login       string    `gorm:"size:64;not null;uniqueIndex" json:"login"`
	Password    string    `gorm:"size:100;not null"            json:"-"`
	IsModerator bool      `gorm:"not null;default:false"       json:"is_moderator"`

	CreatedAt   time.Time `gorm:"<-:create;not null"           json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null"                     json:"updated_at"`
}
package ds

type User struct {
	ID          uint      `gorm:"primaryKey"                   json:"id"`
	Login       string    `gorm:"size:64;not null;uniqueIndex" json:"login"`
	Password    string    `gorm:"size:100;not null"            json:"password"`
	IsModerator bool      `gorm:"not null;default:false"       json:"is_moderator"`
}
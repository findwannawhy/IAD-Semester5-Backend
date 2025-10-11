package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID          uuid.UUID   `json:"id"`
	Login       string `json:"login"`
	IsModerator bool   `json:"is_moderator"`
}
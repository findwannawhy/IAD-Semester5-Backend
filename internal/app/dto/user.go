package dto

import "github.com/google/uuid"

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID          uuid.UUID   `json:"id"`
	Login       string `json:"login"`
	IsModerator bool   `json:"is_moderator"`
}
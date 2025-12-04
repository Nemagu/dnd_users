package webschema

import "github.com/google/uuid"

type UserResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	State  string    `json:"state"`
	Status string    `json:"status"`
}

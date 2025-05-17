package user

import (
	"github.com/google/uuid"
)

type Response struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

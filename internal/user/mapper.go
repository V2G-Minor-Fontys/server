package user

import (
	"github.com/V2G-Minor-Fontys/server/internal/repository"
)

func mapUserToResponse(u *User) *Response {
	return &Response{
		ID:       u.ID,
		Username: u.Username,
	}
}

func mapDatabaseUserToUser(u *repository.User) *User {
	return &User{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}

package repository

import (
	"context"
	"github.com/MrWhok/IMK-FP-BACKEND/entity"
)

type UserRepository interface {
	Authentication(ctx context.Context, username string) (entity.User, error)
	Create(username string, password string, roles []string) error
	DeleteAll()
}

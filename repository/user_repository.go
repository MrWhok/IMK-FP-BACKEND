package repository

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
)

type UserRepository interface {
	Authentication(ctx context.Context, username string) (entity.User, error)
	Create(username string, password string, roles []string, address string, phone string, email string, firstName string, lastName string) error
	DeleteAll()
	FindByUsername(ctx context.Context, username string) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	UpdateProfile(ctx context.Context, username string, email string, phone string, address string) error
	FindAllOrderedByPoints(ctx context.Context) ([]entity.User, error)
}

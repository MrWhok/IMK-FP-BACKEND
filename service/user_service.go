package service

import (
	"context"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
)

type UserService interface {
	Authentication(ctx context.Context, model model.UserModel) entity.User
	Register(ctx context.Context, model model.UserCreateModel) error
	FindMe(ctx context.Context, username string) (model.UserCreateModel, error)
	GetLeaderboard(ctx context.Context) ([]model.UserLeaderboardModel, error)
	UpdateProfile(ctx context.Context, username string, model model.UserUpdateModel) error
}

package service

import (
	"context"
	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	
)

type UserService interface {
	Authentication(ctx context.Context, model model.UserModel) entity.User
	Register(ctx context.Context, model model.UserModel) error
}

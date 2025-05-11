package impl

import (
	"context"
	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

func NewUserServiceImpl(userRepository *repository.UserRepository) service.UserService {
	return &userServiceImpl{UserRepository: *userRepository}
}

type userServiceImpl struct {
	repository.UserRepository
}

func (userService *userServiceImpl) Authentication(ctx context.Context, model model.UserModel) entity.User {
	userResult, err := userService.UserRepository.Authentication(ctx, model.Username)
	if err != nil {
		panic(exception.UnauthorizedError{
			Message: err.Error(),
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(model.Password))
	if err != nil {
		panic(exception.UnauthorizedError{
			Message: "incorrect username and password",
		})
	}
	return userResult
}

func (u *userServiceImpl) Register(ctx context.Context, user model.UserModel) error {

	_,err := u.UserRepository.Authentication(ctx, user.Username)
	if err == nil {
		return fmt.Errorf("username already exists")
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	roles := user.Roles
	if len(roles) == 0 {
		roles = []string{"user"}
	}

	err = u.UserRepository.Create(user.Username, string(hasedPassword), roles)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil

}
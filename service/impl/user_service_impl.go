package impl

import (
	"context"
	"fmt"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/MrWhok/IMK-FP-BACKEND/service"
	"golang.org/x/crypto/bcrypt"
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

func (u *userServiceImpl) Register(ctx context.Context, user model.UserCreateModel) error {

	_, err := u.UserRepository.Authentication(ctx, user.Username)
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

	err = u.UserRepository.Create(user.Username, string(hasedPassword), roles, user.Address, user.Phone, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil

}

func (u *userServiceImpl) FindMe(ctx context.Context, username string) (model.UserCreateModel, error) {
	usernameResult, err := u.UserRepository.Authentication(ctx, username)
	if err != nil {
		return model.UserCreateModel{}, fmt.Errorf("user not found")
	}

	var userRoles []string
	for _, userRole := range usernameResult.UserRoles {
		userRoles = append(userRoles, userRole.Role)
	}

	users, _ := u.UserRepository.FindAllOrderedByPoints(ctx)
	rank := 0
	for i, user := range users {
		if user.Username == username {
			rank = i + 1
			break
		}
	}

	return model.UserCreateModel{
		Username:  usernameResult.Username,
		FirstName: usernameResult.FirstName,
		LastName:  usernameResult.LastName,
		Email:     usernameResult.Email,
		Phone:     usernameResult.Phone,
		Address:   usernameResult.Address,
		Roles:     userRoles,
		Points:    usernameResult.Points,
		Rank:      rank,
	}, nil
}

func (s *userServiceImpl) GetLeaderboard(ctx context.Context) ([]model.UserLeaderboardModel, error) {
	users, err := s.UserRepository.FindAllOrderedByPoints(ctx)
	if err != nil {
		return nil, err
	}

	var leaderboard []model.UserLeaderboardModel
	for i, user := range users {
		leaderboard = append(leaderboard, model.UserLeaderboardModel{
			Rank:      i + 1,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Points:    user.Points,
		})
	}

	return leaderboard, nil
}

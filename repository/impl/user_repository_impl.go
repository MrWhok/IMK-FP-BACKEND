package impl

import (
	"context"
	"errors"

	"github.com/MrWhok/IMK-FP-BACKEND/entity"
	"github.com/MrWhok/IMK-FP-BACKEND/exception"
	"github.com/MrWhok/IMK-FP-BACKEND/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewUserRepositoryImpl(DB *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{DB: DB}
}

type userRepositoryImpl struct {
	*gorm.DB
}

func (userRepository *userRepositoryImpl) Create(username string, password string, roles []string, address string, phone string, email string, firtName string, lastName string) error {
	var userRoles []entity.UserRole
	for _, role := range roles {
		userRoles = append(userRoles, entity.UserRole{
			Id:       uuid.New(),
			Username: username,
			Role:     role,
		})
	}
	user := entity.User{
		Username:  username,
		Password:  password,
		IsActive:  true,
		UserRoles: userRoles,
		FirstName: firtName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Address:   address,
	}

	err := userRepository.DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (userRepository *userRepositoryImpl) DeleteAll() {
	err := userRepository.DB.Where("1=1").Delete(&entity.User{}).Error
	exception.PanicLogging(err)
}

func (userRepository *userRepositoryImpl) Authentication(ctx context.Context, username string) (entity.User, error) {
	var userResult entity.User
	result := userRepository.DB.WithContext(ctx).
		Joins("inner join tb_user_role on tb_user_role.username = tb_user.username").
		Preload("UserRoles").
		Where("tb_user.username = ? and tb_user.is_active = ?", username, true).
		Find(&userResult)
	if result.RowsAffected == 0 {
		return entity.User{}, errors.New("user not found")
	}
	return userResult, nil
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	result := r.DB.WithContext(ctx).Where("username = ?", username).First(&user)
	return user, result.Error
}

func (r *userRepositoryImpl) Update(ctx context.Context, user entity.User) error {
	return r.DB.WithContext(ctx).Save(&user).Error
}

func (r *userRepositoryImpl) FindAllOrderedByPoints(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	result := r.DB.WithContext(ctx).Order("points DESC").Find(&users)
	return users, result.Error
}

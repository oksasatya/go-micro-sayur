package repository

import (
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"
)

type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	modelUser := &model.User{}
	if err := u.db.Where("email = ? && is_verified = ?", email, true).Preload("Roles").First(modelUser).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] GetUserByEmail: user not found with email %s", email)
			return nil, err
		}
		log.Errorf("[UserRepository-1] GetUserByEmail: %v", err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Email:      modelUser.Email,
		Password:   modelUser.Password,
		RoleName:   modelUser.Roles[0].Name,
		Address:    modelUser.Address,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		IsVerified: modelUser.IsVerified,
	}, nil
}

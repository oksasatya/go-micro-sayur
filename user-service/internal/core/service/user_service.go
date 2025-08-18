package service

import (
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/utils/conv"
)

type UserServiceInterface interface {
	SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error)
}

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error) {
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[UserService-1] SignIn: %v", err)
		return nil, "", err
	}
	if checkPass := conv.CheckPasswordHash(req.Password, user.Password); !checkPass {
		err = errors.New("invalid email or password")
		log.Errorf("[UserService-2] SignIn: %v", err)
		return nil, "", err
	}

	return user, "", nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user-service/config"
	"user-service/internal/adapter/message"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/utils/conv"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type UserServiceInterface interface {
	SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error)
	CreateUserAccount(ctx context.Context, req *entity.UserEntity) error
	ForgotPassword(ctx context.Context, req *entity.UserEntity) error
}

type UserService struct {
	repo       repository.UserRepositoryInterface
	cfg        *config.Config
	jwtService JwtServiceInterface
	repoToken  repository.VerificationTokenRepositoryInterface
}

func (u *UserService) ForgotPassword(ctx context.Context, req *entity.UserEntity) error {
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[UserService-1] ForgotPassword: %v", err)
		return err
	}
	token := uuid.New().String()
	reqEntity := &entity.VerificationTokenEntity{
		UserID:    user.ID,
		Token:     token,
		TokenType: "forgot_password",
	}

	err = u.repoToken.CreateVerificationToken(ctx, *reqEntity)
	if err != nil {
		log.Errorf("[UserService-2] ForgotPassword: failed to create verification token: %v", err)
		return err
	}
	urlForgot := fmt.Sprintf("%s/reset-password?token=%s", u.cfg.App.UrlForgotPassword, token)
	messageParam := fmt.Sprintf("Please reset your password by clicking the link: %s/reset-password?token=%s", urlForgot, token)
	err = message.PublishMessage(req.Email, messageParam, "forgot_password")
	if err != nil {
		log.Errorf("[UserService-3] ForgotPassword: failed to publish message: %v", err)
	}
	return nil
}

func (u *UserService) CreateUserAccount(ctx context.Context, req *entity.UserEntity) error {
	password, err := conv.HashPassword(req.Password)
	if err != nil {
		log.Errorf("[UserService-1] CreateUserAccount: failed to hash password: %v", err)
		return err
	}
	req.Password = password
	token := uuid.New().String()
	req.Token = token

	err = u.repo.CreateUserAccount(ctx, req)
	if err != nil {
		log.Errorf("[UserService-2] CreateUserAccount: failed to create user account: %v", err)
		return err
	}

	urlVerify := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)
	messageparam := fmt.Sprintf("Please verify your account by clicking the link: %s/verify?token=%s", urlVerify, req.Token)
	err = message.PublishMessage(req.Email, messageparam, "email_verification")
	if err != nil {
		log.Errorf("[UserService-3] CreateUserAccount: failed to publish message: %v", err)
		return err
	}
	return nil
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
	token, err := u.jwtService.GenerateToken(user.ID)
	if err != nil {
		log.Errorf("[UserService-3] SignIn: failed to generate token: %v", err)
		return nil, "", err
	}

	sessionData := map[string]interface{}{
		"user_id":    user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"logged_in":  true,
		"created_at": time.Now().String(),
		"token":      token,
	}

	redisConn := config.NewRedisClient()
	err = redisConn.HSet(ctx, token, sessionData).Err()
	if err != nil {
		log.Errorf("[UserService-4] SignIn: failed to set session data in Redis: %v", err)
		return nil, "", err
	}

	return user, token, nil
}

func NewUserService(repo repository.UserRepositoryInterface, cfg *config.Config, jwtService JwtServiceInterface, repoToken repository.VerificationTokenRepositoryInterface) UserServiceInterface {
	return &UserService{
		repo:       repo,
		cfg:        cfg,
		jwtService: jwtService,
		repoToken:  repoToken,
	}
}

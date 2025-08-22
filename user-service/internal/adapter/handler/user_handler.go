package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/handler/request"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service"
)

type UserHandlerInterface interface {
	SignIn(ctx echo.Context) error
	CreateUserAccount(ctx echo.Context) error
}

type UserHandler struct {
	userService service.UserServiceInterface
}

var err error

func NewUserHandler(e *echo.Echo, userService service.UserServiceInterface, cfg *config.Config) UserHandlerInterface {
	userHandler := &UserHandler{
		userService: userService,
	}
	e.Use(middleware.Recover())
	e.POST("/signIn", userHandler.SignIn)
	e.POST("/signUp", userHandler.CreateUserAccount)

	// Middleware for checking token
	mid := adapter.NewMiddlewareAdapter(cfg)
	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/check", func(c echo.Context) error {
		return c.String(http.StatusOK, "Admin Check OK")
	})
	return userHandler
}

func (h *UserHandler) SignIn(c echo.Context) error {
	// Implement the SignIn logic here
	var (
		req        = request.SignInRequest{}
		resp       = response.DefaultResponse{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	if err = c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-1] SignIn: failed to bind request: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[UserHandler-1] SignIn: request validation failed: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}
	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}
	user, token, err := h.userService.SignIn(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			log.Infof("[UserHandler-2] SignIn: user not found with email %s", "User Not Found")
			resp.Message = "user not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		log.Errorf("[UserHandler-3] SignIn: failed to sign in: %v", err)
		resp.Message = err.Error()
	}
	respSignIn.ID = user.ID
	respSignIn.Name = user.Name
	respSignIn.Email = user.Email
	respSignIn.Role = user.RoleName
	respSignIn.Lat = user.Lat
	respSignIn.Lng = user.Lng
	respSignIn.Phone = user.Phone
	respSignIn.AccessToken = token

	resp.Message = "success"
	resp.Data = respSignIn
	return c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) CreateUserAccount(c echo.Context) error {
	// Implement the SignUp logic here
	var (
		req  = request.SignUpRequest{}
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)

	if err = c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-1] CreateUserAccount: failed to bind request: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("[UserHandler-1] CreateUserAccount: request validation failed: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if req.Password != req.PasswordConfirmation {
		log.Errorf("[UserHandler-2] CreateUserAccount: password confirmation does not match")
		resp.Message = "password confirmation does not match"
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}
	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err = h.userService.CreateUserAccount(ctx, &reqEntity)
	if err != nil {
		log.Errorf("[UserHandler-3] CreateUserAccount: failed to create user account: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusCreated, resp)
}

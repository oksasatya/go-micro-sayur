package service

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"user-service/config"
)

type JwtServiceInterface interface {
	GenerateToken(userId int64) (string, error)
	ValidateToken(encodeToken string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJwtService(cfg *config.Config) JwtServiceInterface {
	return &jwtService{
		secretKey: cfg.App.JWTSecretKey,
		issuer:    cfg.App.JWTIssuer,
	}
}

// GenerateToken creates a new JWT token for the given user ID.
func (j *jwtService) GenerateToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"iss":     j.issuer,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken checks the validity of the JWT token and returns the parsed token if valid.
func (j *jwtService) ValidateToken(encodeToken string) (*jwt.Token, error) {
	return jwt.Parse(encodeToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	})
}

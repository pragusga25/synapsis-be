package utils

import (
	"errors"
	"synapsis/config"
	"synapsis/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Email  string      `json:"email"`
	UserID uint        `json:"user_id"`
	Role   models.Role `json:"role"`
}

type Claims struct {
	jwt.RegisteredClaims
	User *UserClaims `json:"user"`
}

func GenerateJWT(user *models.User) (string, error) {
	cfg, _ := config.LoadConfig()
	var jwtSecret = []byte(cfg.JWTSecret)

	userRole := models.Role(user.Role)
	if !userRole.IsValid() {
		return "", errors.New("invalid role")
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
		User: &UserClaims{
			Email:  user.Email,
			UserID: user.ID,
			Role:   user.Role,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)

}

func ValidateJWT(tokenString string) (*UserClaims, error) {
	cfg, _ := config.LoadConfig()
	var jwtSecret = []byte(cfg.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := claims["user"].(map[string]interface{})
		return &UserClaims{
			Email:  user["email"].(string),
			UserID: uint(user["user_id"].(float64)),
			Role:   models.Role(user["role"].(string)),
		}, nil
	}

	return nil, errors.New("invalid token")

}

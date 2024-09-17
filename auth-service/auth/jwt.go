package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var (
	jwtSecret          = []byte(os.Getenv("JWT_SECRET")) // Ensure JWT_SECRET is set
	accessTokenExpiry  = 15 * time.Minute                // Short-lived token
	refreshTokenExpiry = 7 * 24 * time.Hour              // Longer-lived token
)

type Claims struct {
	UserID  uint   `json:"user_id"`
	Role    string `json:"role"`
	Version string `json:"version"`
	jwt.StandardClaims
}

// GenerateAccessToken creates a new access token
func GenerateAccessToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExpiry).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken creates a new refresh token
func GenerateRefreshToken(userID uint) (string, error) {
	claims := Claims{
		UserID:  userID,
		Role:    "user",
		Version: "REFRESH",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenExpiry).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

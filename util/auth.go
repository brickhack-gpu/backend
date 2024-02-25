package util

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	"gpu/model"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(username string, userid int64, admin bool, jwtSecret string) (string, string, error) {
	jwtSecretBytes := []byte(jwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      userid,
		"exp":      time.Now().Add(48 * time.Hour).Unix(), // TODO: Refresh tokens
		"username": username,
		"admin":    admin,
	})

	tokenString, err := token.SignedString(jwtSecretBytes)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userid,
		"exp": time.Now().Add(96 * time.Hour).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecretBytes)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func GenerateJWTFromRefreshToken(db *bun.DB, jwtSecret string, rt string, ctx context.Context) (string, error) {
	token, err := jwt.Parse(rt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := new(model.User)
		err := db.NewSelect().Model(&user).Where("id = ?", claims["sub"].(string)).Scan(ctx)
		if err != nil {
			return "", err
		}

		if user.Active {
			jwtSecretBytes := []byte(jwtSecret)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":      user.ID,
				"exp":      time.Now().Add(48 * time.Minute).Unix(), // TODO: Refresh tokens
				"username": user.Username,
				"admin":    user.Admin,
			})

			tokenString, err := token.SignedString(jwtSecretBytes)
			if err != nil {
				return "", err
			}

			return tokenString, err
		} else {
			return "", fmt.Errorf("Disabled")
		}
	}

	return "", err
}

package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/uptrace/bun"

	"gpu/util"
)

type Router struct {
	DB            *bun.DB
	JwtSecret     string
	StripeSecret  string
	StripeWebhook string
	Dev           bool
}

func NewRouter(db *bun.DB, jwtSecret, stripeSecret, stripeWebhook string, dev bool) *Router {
	return &Router{
		DB:            db,
		JwtSecret:     jwtSecret,
		StripeSecret:  stripeSecret,
		StripeWebhook: stripeWebhook,
		Dev:           dev,
	}
}

func (router *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			util.ResError(errors.New(""), w, http.StatusUnauthorized, "Malformed token")
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(router.JwtSecret), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "props", claims)
				// Access context values in handlers like this
				// props, _ := r.Context().Value("props").(jwt.MapClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				util.ResError(err, w, http.StatusUnauthorized, "Unauthorized")
			}
		}
	})
}

package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/farhanfatur/simple-example-graphql-golang/internal/pkg/jwt"
	"github.com/farhanfatur/simple-example-graphql-golang/internal/users"
)

var userContextKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Check authenticated user
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check validate token
			username, err := jwt.ParseToken(header)
			if err != nil {
				http.Error(w, "Invalid token!!", http.StatusForbidden)
				return
			}

			// check if username is exist in db
			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(user.Username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user.ID = strconv.Itoa(id)

			ctx := context.WithValue(r.Context(), userContextKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userContextKey).(*users.User)
	return raw
}

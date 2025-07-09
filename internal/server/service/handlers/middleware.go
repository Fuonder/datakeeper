package handlers

import (
	"context"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	"net/http"
	"time"
)

func (h Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return logger.HanlderWithLogger(func(rw http.ResponseWriter, r *http.Request) {
		logger.Log.Debug("Auth middleware")

		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie == nil {
			SendResponse(rw, http.StatusUnauthorized, []byte("Missing or invalid token"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tokenString := cookie.Value
		err = h.authSrv.ValidateJWT(ctx, tokenString)
		if err != nil {
			SendResponse(rw, http.StatusUnauthorized, []byte("Invalid token"))
			return
		}

		next.ServeHTTP(rw, r)
	})
}

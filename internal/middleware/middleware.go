package middleware

import (
	"context"
	"net/http"
	"strconv"
	"vk-intern_test-case/utils/database"
	"vk-intern_test-case/utils/response"

	log "github.com/sirupsen/logrus"
)

const (
	adminCheck = `select role from service_user where id = $1`
)

type AuthMiddleware struct {
	pool database.PgxIface
}

func NewAuthMiddleware(pool database.PgxIface) *AuthMiddleware {
	return &AuthMiddleware{
		pool: pool,
	}
}

func (aM *AuthMiddleware) MiddlewareCheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Middleware started")
		if r.Method != http.MethodGet {
			jsonEnc := response.MakeJsonEncoder(w)
			headerValue, exist := r.Header[http.CanonicalHeaderKey("Authorization")]
			log.Debug(headerValue)
			if !exist {
				response.WriteBasicResponse(w, jsonEnc, http.StatusUnauthorized, "Unauthorized")
				return
			}

			userID, err := strconv.Atoi(headerValue[0])
			if err != nil {
				response.WriteBasicResponse(w, jsonEnc, http.StatusUnauthorized, "Unauthorized")
			}

			tx, err := aM.pool.Begin(context.Background())
			if err != nil {
				response.WriteBasicResponse(w, jsonEnc, http.StatusUnauthorized, "Unauthorized")
			}
			var role string
			row := tx.QueryRow(context.Background(), adminCheck, &userID)
			err = row.Scan(&role)
			if err != nil {
				response.WriteBasicResponse(w, jsonEnc, http.StatusUnauthorized, "Unauthorized")
			}

			if role != "Администратор" {
				response.WriteBasicResponse(w, jsonEnc, http.StatusUnauthorized, "Unauthorized")
			}
		}

		next.ServeHTTP(w, r)
	})
}
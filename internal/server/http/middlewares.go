package http

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := negroni.NewResponseWriter(w)

		start := time.Now()

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		h.logger.Info("Request info HTTP",
			zap.String("client ip", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("method path", r.URL.Path),
			zap.String("HTTP version", r.Proto),
			zap.String("status code", strconv.Itoa(lrw.Status())),
			zap.String("processing time", duration.String()),
		)
	})
}

func (h *Handler) authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return nil, errors.New("error parsing token")
			}
			return []byte(h.secretKey), nil
		})
		if err != nil || !token.Valid {
			resp := newResponse("", "unauthorized", err)
			h.logError(resp.Message, getUserAction, resp.Error)
			renderResponse(w, r, http.StatusUnauthorized, resp)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := claims["user_id"].(string)
			ctx := context.WithValue(r.Context(), "user_id", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

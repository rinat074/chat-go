package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rinat074/chat-go/pkg/logger"
	"go.uber.org/zap"
)

// Recovery middleware для восстановления после паники в HTTP обработчиках
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Логирование паники и стек-трейса
				stack := debug.Stack()
				logger.Error("Recovered from panic",
					zap.Any("error", err),
					zap.String("stack", string(stack)),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("remote_addr", r.RemoteAddr),
				)

				// Отправка 500 Internal Server Error клиенту
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"Internal Server Error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// GRPCRecovery возвращает функцию восстановления для gRPC сервера
func GRPCRecovery() func(interface{}) error {
	return func(p interface{}) error {
		stack := debug.Stack()
		logger.Error("Recovered from panic in gRPC handler",
			zap.Any("error", p),
			zap.String("stack", string(stack)),
		)
		return fmt.Errorf("internal error")
	}
}

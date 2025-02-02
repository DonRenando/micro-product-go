package middleware

import (
	"net/http"

	"github.com/shadowshot-x/micro-product-go/authservice/jwt"
	"go.uber.org/zap"
)

// TokenMiddleware is the token validation route handler
type TokenMiddleware struct {
	logger *zap.Logger
}

// NewTokenMiddleware returns a frsh Token controller
func NewTokenMiddleware(logger *zap.Logger) *TokenMiddleware {
	return &TokenMiddleware{
		logger: logger,
	}
}

// Middleware itself returns a function that is a Handler. it is executed for each request.
// We want all our routes for REST to be authenticated. So, we validate the token
func (ctrl *TokenMiddleware) TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// check if token is present
		if _, ok := r.Header["Token"]; !ok {
			ctrl.logger.Warn("Token was not found in the header")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Missing"))
			return
		}
		token := r.Header["Token"][0]
		check, err := jwt.ValidateToken(token, "S0m3_R4n90m_sss")
		if err != nil {
			ctrl.logger.Error("Token Validation Failed", zap.String("token", token))
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Token Validation Failed"))
			return
		}
		if !check {
			ctrl.logger.Warn("Token invalid", zap.String("token", token))
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Invalid"))
			return
		}
		// rw.WriteHeader(http.StatusOK)
		// rw.Write([]byte("Authorized Token"))

		// this calls the next function. If not included, the router wont entertain any requests
		next.ServeHTTP(rw, r)
	})
}

package auth

import (
	"net/http"
)

type AuthService interface {
	ContextInitiator(next http.Handler) http.Handler
	Encode(claims map[string]interface{}) (string, error)
	Decode(jwt string) (map[string]interface{}, error)

	Authenticator(tokenFn func(r *http.Request) (string, error)) (mw func(http.Handler) http.Handler)
	GetTokenFromHeader(r *http.Request) (string, error)
	GetTokenFromCookie(r *http.Request) (string, error)
}

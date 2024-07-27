package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"
	"github.com/golang-jwt/jwt"
)

type authService struct {
	tokenKey   string
	expiryTime int // time in hours
	secretKey  string
	log        logger.Logger
}

func NewAuthService(signingKey string, expiryTime int, tokenKey string, log logger.Logger) AuthService {
	return &authService{
		secretKey:  signingKey,
		expiryTime: expiryTime,
		tokenKey:   tokenKey,
		log:        log,
	}
}

func (as *authService) ContextInitiator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		ctx := context.WithValue(r.Context(), model.TokenKey, token)
		req := r.Clone(ctx)
		next.ServeHTTP(w, req)
	})
}

func (as *authService) Decode(rawToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(as.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (as *authService) Encode(claims map[string]interface{}) (string, error) {
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(as.expiryTime)).Unix()
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(as.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (as *authService) Authenticator(tokenFn func(r *http.Request) (string, error)) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken, err := tokenFn(r)
			if err != nil {
				as.log.Errorf("error getting user from request %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Get the unique request ID for this request and add it to the request context
			// This is also referred to as the "correlation ID"
			reqID, err := model.GetReqIDFromRequest(r)
			if err != nil {
				// we just warn and move fordward
				// all request should contain a request ID in the future
				as.log.
					Warnf(err.Error())
				reqID = "not-set"
			}
			ctx := context.WithValue(r.Context(), model.ReqIDKey, reqID)

			claims, err := as.Decode(rawToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			user, err := model.ClaimsToUser(claims)
			if err != nil {
				as.log.Errorf("failed to get user from e_commerce_webapp_token: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := user.Validate(); err != nil {
				as.log.Errorf("failed to validate user claims: %v", err)
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			as.log.
				WithReqID(ctx).
				Infof("user token: %v", rawToken)

			ctx = context.WithValue(ctx, model.UserKey, user)
			req := r.Clone(ctx)

			as.log.
				WithReqID(ctx).
				Infof("request from user %v to path %v, user-groups: %v", user.Username, r.URL.Path, user.FederatedGroups)
			h.ServeHTTP(w, req)
		})
	}
	return
}

func (as *authService) GetTokenFromHeader(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errors.New("token not found")
	}
	token = strings.TrimPrefix(token, "Bearer ")
	return token, nil
}

func (as *authService) GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(as.tokenKey)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

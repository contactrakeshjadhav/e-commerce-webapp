package middleware

import (
	"context"
	"net/http"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/logger"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
	"github.com/pkg/errors"
)

var ErrNoAllowedActionsForProject error = errors.New("cant find allowed actions for user on given project")

type authorizationMiddleware struct {
	log logger.Logger
}

func NewAuthorizationMiddleware(
	log logger.Logger,
) *authorizationMiddleware {
	return &authorizationMiddleware{

		log: log,
	}
}

// ValidateAdminAccess Will validate if root project has admin access.
// 403 status code will be returned if user has access denied.
func (a *authorizationMiddleware) ValidateAdminAccess(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, ok := commons.GetUserFromCTX(ctx)
		if !ok {
			a.log.Errorf("validate admin access - failed to get user from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

	})
}

// isAdmin Returns true if user is admin in the project
func (a *authorizationMiddleware) isAdmin(ctx context.Context, projectID model.ID) (bool, error) {
	return false, nil
}

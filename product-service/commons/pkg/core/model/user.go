package model

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// User is the object that define the properties of the OIDC claims
type User struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Email           string   `json:"email,omitempty"`
	Username        string   `json:"username,omitempty"`
	FederatedID     string   `json:"federatedId,omitempty"`
	FederatedGroups []string `json:"federatedGroups,omitempty"`
	LastLogin       string   `json:"lastLogin,omitempty"`

	// AccessedWorkspaces is a set of Workspace IDs.
	// ID is added when user connects to the Workspace.
	AccessedWorkspaces map[string]interface{} `json:"accessed_workspaces,omitempty"`
}

// NewUser creates a new User object instance
func NewUser(name string, email string, username string, federatedId string, federatedGroups []string) User {
	return User{
		Name:               name,
		Email:              email,
		Username:           username,
		FederatedID:        federatedId,
		FederatedGroups:    federatedGroups,
		AccessedWorkspaces: map[string]interface{}{},
	}
}

// Validate ensure that the user has the ID, username, name and email properties set
func (u *User) Validate() error {
	if len(strings.TrimSpace(u.FederatedID)) == 0 {
		return errors.Errorf("user Federated ID is required")
	}

	if len(strings.TrimSpace(u.Username)) == 0 {
		return errors.Errorf("user USERNAME is required")
	}

	if len(strings.TrimSpace(u.Email)) == 0 {
		return errors.Errorf("user EMAIL is required")
	}

	if len(strings.TrimSpace(u.Name)) == 0 {
		return errors.Errorf("user NAME is required")
	}
	return nil
}

func (u *User) ToClaims() map[string]interface{} {
	return map[string]interface{}{
		"id":                  u.ID,
		"name":                u.Name,
		"uid":                 u.FederatedID,
		"roles":               u.FederatedGroups,
		"email":               u.Email,
		"username":            u.Username,
		"accessed_workspaces": u.AccessedWorkspaces,
	}
}

var ErrClaimHasInvalidType = errors.New("claim has invalid type")

func ClaimsToUser(claims map[string]interface{}) (User, error) {
	idValue, ok := claims["id"]
	if !ok {
		return User{}, errors.New("no user id found in token")
	}
	id := idValue.(string)

	uidValue, ok := claims["uid"]
	if !ok {
		return User{}, errors.New("no user uid found in token")
	}
	uid := uidValue.(string)

	usernameValue, ok := claims["username"]
	if !ok {
		return User{}, errors.New("no user username found in token")
	}
	username := usernameValue.(string)

	emailValue, ok := claims["email"]
	if !ok {
		return User{}, errors.New("no user email found in token")
	}
	email := emailValue.(string)

	nameValue, ok := claims["name"]
	if !ok {
		return User{}, errors.New("no user name found in token")
	}
	name := nameValue.(string)

	rolesValue, ok := claims["roles"]
	var roles []string
	if !ok {
		return User{}, errors.New("no user roles found in token")
	}
	rolesInterface, ok := rolesValue.([]interface{})
	if ok {
		roles = make([]string, len(rolesInterface))
		for i, role := range rolesInterface {
			roles[i] = fmt.Sprint(role)
		}
	} else {
		if rolesValue != nil {
			return User{},
				errors.Wrapf(
					ErrClaimHasInvalidType,
					"list of roles is expected to be of type %v but %v was given",
					reflect.TypeOf([]interface{}{}),
					reflect.TypeOf(rolesValue),
				)
		}
		roles = []string{}
	}

	var accessedWorkspaces map[string]interface{}
	accessedWorkspacesValue, ok := claims["accessed_workspaces"]
	if ok {
		accessedWorkspaces, ok = accessedWorkspacesValue.(map[string]interface{})
		if !ok {
			if accessedWorkspacesValue != nil {
				return User{}, errors.Wrapf(ErrClaimHasInvalidType,
					"set of accessed projects is expected to be of type %v but %v was given",
					reflect.TypeOf(map[string]interface{}{}),
					reflect.TypeOf(accessedWorkspacesValue),
				)
			}
			accessedWorkspaces = map[string]interface{}{}
		}
	}

	return User{
		ID:                 id,
		FederatedID:        uid,
		Username:           username,
		Email:              email,
		FederatedGroups:    roles,
		Name:               name,
		AccessedWorkspaces: accessedWorkspaces,
	}, nil
}

func GetUserFromCTX(ctx context.Context) (User, bool) {
	rawUser := ctx.Value(UserKey)
	if rawUser == nil {
		return User{}, false
	}

	user, ok := rawUser.(User)
	if !ok {
		return User{}, false
	}

	return user, true
}

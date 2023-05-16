package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/phantranhieunhan/authen-author/db/sqlc"
	"github.com/phantranhieunhan/authen-author/token"
	"github.com/phantranhieunhan/authen-author/util"
)

var (
	errorGetUser         = errors.New("errorGetUser")
	sttNoDeciderProvided = errors.New("stt no decider provided")
	sttDeniedAll         = errors.New("sttDeniedAll")
	sttNotAllowed        = errors.New("sttNotAllowed")
	userNoProvided       = errors.New("userNoProvided")
)

type Url struct {
	Method string
	Path   string
}

var rbacDecider = map[Url][]string{
	{Method: http.MethodGet, Path: "/accounts/:id"}:  {RoleCentreStaff},
	{Method: http.MethodPost, Path: "/accounts/:id"}: {RoleSchoolAdmin, RoleCentreLead},
	{Method: http.MethodGet, Path: "/demo"}: nil, // no need permission
}

type GroupDecider struct {
	GroupFetcher  func(ctx context.Context, userID string) ([]string, error)
	AllowedGroups map[Url][]string
}

func NewGroupDecider(store db.Store) *GroupDecider {
	return &GroupDecider{
		GroupFetcher: func(ctx context.Context, username string) ([]string, error) {
			user, err := store.GetUser(ctx, username)
			if err != nil {
				if err == sql.ErrNoRows {
					return []string{}, nil
				}
				return nil, errorGetUser
			}
			return user.Role, nil
		},
		AllowedGroups: rbacDecider,
	}
}

// Check checks if user allowed to call a method
func (g *GroupDecider) Check(ctx context.Context, userID, method, path string) (groups []string, err error) {
	allowedGroups, ok := g.AllowedGroups[Url{Method: method, Path: path}]
	if !ok {
		return nil, sttNoDeciderProvided
	}
	groups, _ = g.GroupFetcher(ctx, userID)
	if len(groups) == 0 {
		return nil, sttDeniedAll
	}

	if allowedGroups == nil {
		// allowed all
		return groups, nil
	}

	if len(allowedGroups) == 0 {
		return nil, sttDeniedAll
	}

	for _, group := range groups {
		if util.InArrayString(group, allowedGroups) {
			return groups, nil
		}
	}
	return nil, sttNotAllowed
}

func authorMiddleware(store db.Store) gin.HandlerFunc {
	gd := NewGroupDecider(store)

	return func(ctx *gin.Context) {
		authPayload, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		if !ok || authPayload == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(userNoProvided))
			return
		}
		// ctx.Pa
		_, err := gd.Check(ctx, authPayload.Username, ctx.Request.Method, ctx.FullPath())
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.Next()
	}
}

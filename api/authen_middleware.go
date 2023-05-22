package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/authen-author/db/redis"
	db "github.com/phantranhieunhan/authen-author/db/sqlc"
	"github.com/phantranhieunhan/authen-author/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authenMiddleware(tokenMaker token.Maker, store db.Store, session redis.SessionStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New("authorization type is not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// check session is available on DB
		count, err := store.CheckIsAvailable(ctx, payload.SourceID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if count == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("Session is expired."))
			return
		}

		// get token in redis
		// token, err := session.GetToken(ctx, fmt.Sprintf(AccessTokenIDPrefix, payload.ID.String()))
		// if err != nil && err != redis.ErrNilSession {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }
		// if token != "" {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("Token is expired.")))
		// 	return
		// }

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

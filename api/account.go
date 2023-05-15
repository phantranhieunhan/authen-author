package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

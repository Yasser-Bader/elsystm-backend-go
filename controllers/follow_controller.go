package controllers

import (
	"net/http"
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/repositories"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
)

type FollowController struct {
	follow repositories.FollowRepository
}

func NewFollowController() *FollowController {
	return &FollowController{
		follow: repositories.NewFollowRepository(),
	}
}

func (controller *FollowController) Followers(ctx *gin.Context) {
	var pagination util.Pagination
	page, _ := strconv.Atoi(ctx.Query("page"))
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	controller.follow.Followers(page, userId, &pagination)
	ctx.JSON(http.StatusOK, util.GenerateResponse(pagination, true, ""))
}

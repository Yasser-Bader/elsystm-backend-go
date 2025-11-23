package controllers

import (
	"net/http"
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/repositories"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
)

type LikesController struct {
	like repositories.LikeRepository
}

func NewLikesController() *LikesController {
	return &LikesController{
		like: repositories.NewLikeRepository(),
	}
}

func (controller *LikesController) Store(ctx *gin.Context) {

	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	var likeDTO dto.LikeCreateDTO
	if err := ctx.ShouldBind(&likeDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": util.Validation(err)})
		return
	}

	// Store like
	like, err := controller.like.Store(userId, likeDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}

	ctx.JSON(http.StatusOK, util.GenerateResponse(like, true, ""))
}

func (controller *LikesController) Delete(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := controller.like.Delete(userId, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse("Deleted", true, ""))
}

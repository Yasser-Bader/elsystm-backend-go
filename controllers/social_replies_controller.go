package controllers

import (
	"net/http"
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/repositories"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
)

type RepliesController struct {
	reply repositories.ReplyRepository
	media repositories.SocialMediaRepository
}

func NewRepliesController() *RepliesController {
	return &RepliesController{
		reply: repositories.NewReplyRepository(),
		media: repositories.NewSocialMediaRepository(),
	}
}

func (controller *RepliesController) Store(ctx *gin.Context) {

	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	var replyDTO dto.ReplyCreateDTO
	if err := ctx.ShouldBind(&replyDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": util.Validation(err)})
		return
	}

	// Store reply
	reply, err := controller.reply.Store(userId, replyDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}

	// Store Media
	media, err := controller.media.StoreSingleFile(userId, reply.ID, "reply", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	reply.Media = media

	ctx.JSON(http.StatusOK, util.GenerateResponse(reply, true, ""))
}

func (controller *RepliesController) Delete(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := controller.reply.Delete(userId, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse("Deleted", true, ""))
}

func (controller *RepliesController) Update(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	var replyDTO dto.ReplyUpdateDTO
	if err := ctx.ShouldBind(&replyDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": util.Validation(err)})
		return
	}
	replyUpdated, err := controller.reply.Update(userId, id, replyDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse(replyUpdated, true, "Updated"))
}

package controllers

import (
	"net/http"
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/repositories"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
)

type CommentsController struct {
	comment repositories.CommentRepository
	media   repositories.SocialMediaRepository
}

func NewCommentsController() *CommentsController {
	return &CommentsController{
		comment: repositories.NewCommentRepository(),
		media:   repositories.NewSocialMediaRepository(),
	}
}

func (controller *CommentsController) Store(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	var commentDTO dto.CommentCreateDTO
	if err := ctx.ShouldBind(&commentDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": util.Validation(err)})
		return
	}

	// Store comment
	comment, err := controller.comment.Store(userId, commentDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	// Store Media
	socialMedia, err := controller.media.StoreSingleFile(userId, comment.ID, "comment", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Internal error"})
	}

	comment.Media = socialMedia

	ctx.JSON(http.StatusCreated, util.GenerateResponse(comment, true, ""))
}

func (controller *CommentsController) Delete(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := controller.comment.Delete(userId, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse("Deleted", true, ""))
}

func (controller *CommentsController) Update(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	var commentDTO dto.CommentUpdateDTO
	if err := ctx.ShouldBind(&commentDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": util.Validation(err)})
		return
	}
	commentUpdated, err := controller.comment.Update(userId, id, commentDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse(commentUpdated, true, "Updated"))
}

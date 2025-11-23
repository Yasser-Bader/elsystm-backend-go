package controllers

import (
	"net/http"
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/repositories"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
)

type HashtagsController struct {
	hashtag repositories.HashtagRepository
}

func NewHashtagsController() *HashtagsController {
	return &HashtagsController{
		hashtag: repositories.NewHashtagRepository(),
	}
}

func (controller *HashtagsController) Store(ctx *gin.Context) {

	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	var hashtagDTO dto.HashtagCreateDTO
	if err := ctx.ShouldBind(&hashtagDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": util.Validation(err)})
		return
	}
	isDuplicateName := controller.hashtag.HashtagIsDuplicateName(hashtagDTO.Name)

	if isDuplicateName.RowsAffected > 0 {
		util.ValidationReturnErrorResponse(ctx, "name", "unique", "name is available")
		return
	}

	// Store hashtag
	hashtag, err := controller.hashtag.Store(userId, hashtagDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}

	ctx.JSON(http.StatusOK, util.GenerateResponse(hashtag, true, ""))
}

func (controller *HashtagsController) Delete(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := controller.hashtag.Delete(userId, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse("Deleted", true, ""))
}

package controllers

import (
	"net/http"
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/repositories"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
)

type PostsController struct {
	post    repositories.PostRepository
	mention repositories.MentionRepository
	media   repositories.SocialMediaRepository
	hashtag repositories.HashtagRepository
}

func NewPostsController() *PostsController {
	return &PostsController{
		post:    repositories.NewPostRepository(),
		media:   repositories.NewSocialMediaRepository(),
		mention: repositories.NewMentionRepository(),
		hashtag: repositories.NewHashtagRepository(),
	}
}

func (controller *PostsController) Index(ctx *gin.Context) {
	var pagination util.Pagination
	page, _ := strconv.Atoi(ctx.Query("page"))
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	if err := controller.post.Index(userId, page, &pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusOK, util.GenerateResponse(pagination, true, ""))
}

func (controller *PostsController) Store(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	var postDTO dto.PostCreateDTO
	if err := ctx.ShouldBind(&postDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": util.Validation(err)})
		return
	}

	// Store Post
	post, err := controller.post.Store(userId, postDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}

	// Store Mentions
	if len(postDTO.Mentions) > 0 {
		err := controller.mention.Store(userId, post.ID, postDTO.Mentions, "post")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
			return
		}
	}

	// Store Hashtags
	if len(postDTO.Hashtags) > 0 {
		err := controller.hashtag.StoreByModel(userId, post.ID, postDTO.Hashtags, "post")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
			return
		}
	}

	// Store Media
	media, err := controller.media.Store(userId, post.ID, "post", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	post.Media = media

	// Send Notification
	util.WorkerMakeRequest(struct {
		Name    string
		Payload interface{}
	}{
		Name: "post created",
		Payload: struct {
			UserId  uint64
			Content string
		}{
			UserId:  userId,
			Content: post.Content,
		},
	})

	ctx.JSON(http.StatusOK, util.GenerateResponse(post, true, ""))
}

func (controller *PostsController) Delete(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := controller.post.Delete(userId, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse("Deleted", true, ""))
}

func (controller *PostsController) Update(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	var postDTO dto.PostUpdateDTO
	if err := ctx.ShouldBind(&postDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": util.Validation(err)})
		return
	}

	//Store Post
	postUpdate, err := controller.post.Update(userId, id, postDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}

	// Store Media
	mediaUpdated, err := controller.media.StoreUpdate(userId, postUpdate.ID, "post", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error file"})
		return
	}
	postUpdate.Media = mediaUpdated

	ctx.JSON(http.StatusCreated, util.GenerateResponse(mediaUpdated, true, "Updated"))
}

func (controller *PostsController) GetPost(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Request.Header.Get("id"), 10, 64)
	post, err := controller.post.GetPost(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusOK, util.GenerateResponse(post, true, ""))
}

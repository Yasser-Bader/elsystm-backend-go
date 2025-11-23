package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/repositories"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
)

type StoriesController struct {
	story     repositories.SocialStoryRepository
	storyUser repositories.SocialStoryUserRepository
	follow    repositories.FollowRepository
}

func NewStoriesController() *StoriesController {
	return &StoriesController{
		story:     repositories.NewStoryRepository(),
		storyUser: repositories.NewSocialStoryUserRepository(),
		follow:    repositories.NewFollowRepository(),
	}
}

func (controller *StoriesController) Store(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	var storyDTO dto.StoryCreateDTO

	err := json.Unmarshal([]byte(fmt.Sprintf(`{"data":%s}`, ctx.Request.FormValue("data"))), &storyDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Internal error"})
		return
	}

	// Store Story
	stories, err := controller.story.Store(userId, storyDTO, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Internal error"})
		return
	}

	ctx.JSON(http.StatusCreated, util.GenerateResponse(stories, true, ""))
}

func (controller *StoriesController) Delete(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := controller.story.Delete(userId, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	if err := controller.storyUser.Delete(userId, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse("Deleted", true, ""))
}

func (controller *StoriesController) Seen(ctx *gin.Context) {
	var storyUserSeen dto.StoryUserSeen
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)

	if err := ctx.ShouldBindJSON(&storyUserSeen); err != nil {
		ctx.JSON(http.StatusBadRequest, util.Validation(err))
	}

	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := controller.storyUser.Store(userId, storyUserSeen, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusAccepted, util.GenerateResponse("seen", true, ""))
}

func (controller *StoriesController) User(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	stories, err := controller.story.User(id, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusOK, util.GenerateResponse(stories, true, ""))
}

func (controller *StoriesController) Share(ctx *gin.Context) {
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	story, err := controller.story.Share(userId, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	if err := controller.storyUser.Delete(userId, story.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusCreated, util.GenerateResponse(story, true, ""))
}

func (controller *StoriesController) Archive(ctx *gin.Context) {
	var pagination util.Pagination
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	page, _ := strconv.Atoi(ctx.Query("page"))
	if err := controller.story.Archive(page, userId, &pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusOK, util.GenerateResponse(pagination, true, ""))
}

func (controller *StoriesController) UsersSeenStory(ctx *gin.Context) {
	var pagination util.Pagination
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	page, _ := strconv.Atoi(ctx.Query("page"))
	if err := controller.storyUser.UsersSeenStory(userId, page, id, &pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}
	ctx.JSON(http.StatusOK, util.GenerateResponse(pagination, true, ""))
}

func (controller StoriesController) NewStories(ctx *gin.Context) {
	var pagination util.Pagination
	page, _ := strconv.Atoi(ctx.Query("page"))
	userId, _ := strconv.ParseUint(ctx.Request.Header.Get("user"), 10, 64)

	followingId, err := controller.follow.FollowingID(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}

	if err := controller.story.NewStories(userId, followingId, page, &pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}

	ctx.JSON(http.StatusOK, util.GenerateResponse(pagination, true, ""))
}

func (controller *StoriesController) Count(ctx *gin.Context) {
	var canCreate bool = true
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	count, err := controller.story.CountTodayStory(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Error"})
		return
	}

	if count > 15 {
		canCreate = false
	}

	ctx.JSON(http.StatusOK, util.GenerateResponse(canCreate, true, ""))
}

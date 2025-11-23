package v1

import (
	"github.com/Elsystm-Inc/systm-go-social/controllers"
	"github.com/gin-gonic/gin"
)

// var db *gorm.DB = config.ConnectDB()
var controller = controllers.NewStoriesController()

func SetupStoryRoute(route *gin.Engine) {
	social := route.Group("/social/v1")
	{
		api := social.Group("stories")
		{
			api.POST("/", controller.Store)
			api.PUT("/update-seen", controller.Seen)
			api.GET("/archive", controller.Archive)
			api.GET("/:id/count", controller.Count)
			api.GET("/new-stories", controller.NewStories)
			api.GET("/:id/user", controller.User)
			api.GET("/:id/user-seen-story", controller.UsersSeenStory)
			api.PUT("/:id/share", controller.Share)
			api.DELETE("/:id/delete", controller.Delete)
		}
	}
}

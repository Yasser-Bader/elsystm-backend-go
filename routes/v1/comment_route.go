package v1

import (
	"github.com/Elsystm-Inc/systm-go-social/controllers"
	"github.com/gin-gonic/gin"
)

var commentsController = controllers.NewCommentsController()

func SetupCommentRoute(route *gin.Engine) {

	social := route.Group("/social/v1")
	{
		api := social.Group("comments")
		{
			api.POST("/", commentsController.Store)
			api.DELETE("/:id/delete", commentsController.Delete)
			api.PUT("/:id/update", commentsController.Update)
		}

	}
}

package v1

import (
	"github.com/Elsystm-Inc/systm-go-social/controllers"
	"github.com/gin-gonic/gin"
)

var repliesController = controllers.NewRepliesController()

func SetupReplyRoute(route *gin.Engine) {

	social := route.Group("/social/v1")
	{
		api := social.Group("replies")
		{
			api.POST("/", repliesController.Store)
			api.DELETE("/:id/delete", repliesController.Delete)
			api.PUT("/:id/update", repliesController.Update)

		}

	}
}

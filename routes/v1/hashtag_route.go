package v1

import (
	"github.com/Elsystm-Inc/systm-go-social/controllers"
	"github.com/gin-gonic/gin"
)

var hashtagsController = controllers.NewHashtagsController()

func SetupHashtagRoute(route *gin.Engine) {

	social := route.Group("/social/v1")
	{
		api := social.Group("hashtags")
		{
			api.POST("/", hashtagsController.Store)
			api.DELETE("/:id/delete", hashtagsController.Delete)

		}

	}
}

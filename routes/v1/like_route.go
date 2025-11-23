package v1

import (
	"github.com/Elsystm-Inc/systm-go-social/controllers"
	"github.com/gin-gonic/gin"
)

var likesController = controllers.NewLikesController()

func SetupLikeRoute(route *gin.Engine) {

	social := route.Group("/social/v1")
	{
		api := social.Group("likes")
		{
			api.POST("/", likesController.Store)
			api.DELETE("/:id/delete", likesController.Delete)

		}

	}
}

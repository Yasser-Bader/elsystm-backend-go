package v1

import (
	"github.com/Elsystm-Inc/systm-go-social/controllers"
	"github.com/gin-gonic/gin"
)

var mentionsController = controllers.NewMentionsController()

func SetupMentionRoute(route *gin.Engine) {

	social := route.Group("/social/v1")
	{
		api := social.Group("mentions")
		{
			api.DELETE("/:id/delete", mentionsController.Delete)

		}

	}
}

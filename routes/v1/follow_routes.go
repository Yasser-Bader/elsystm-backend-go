package v1

import (
	"github.com/Elsystm-Inc/systm-go-social/controllers"
	"github.com/gin-gonic/gin"
)

var followController = controllers.NewFollowController()

func SetupFollowRoute(route *gin.Engine) {

	social := route.Group("/social/v1")
	{
		api := social.Group("follow")
		{
			api.GET("/followers", followController.Followers)
		}

	}
}

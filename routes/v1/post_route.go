package v1

import (
	"github.com/Elsystm-Inc/systm-go-social/controllers"
	"github.com/gin-gonic/gin"
)

var postsController = controllers.NewPostsController()

func SetupPostRoute(route *gin.Engine) {

	social := route.Group("/social/v1")
	{
		api := social.Group("posts")
		{
			api.GET("/", postsController.Index)
			api.POST("/", postsController.Store)
			api.DELETE("/:id/delete", postsController.Delete)
			api.PUT("/:id/update", postsController.Update)
			// api.GET("/", postsController.GetPost)

		}

	}
}

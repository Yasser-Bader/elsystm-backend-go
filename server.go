package main

import (
	"os"

	"github.com/Elsystm-Inc/systm-go-social/middleware"
	v1 "github.com/Elsystm-Inc/systm-go-social/routes/v1"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	server := gin.Default()

	// Enable Cors
	server.Use(middleware.CORSMiddleware())

	// Social Route
	v1.SetupStoryRoute(server)

	// Follow Route
	v1.SetupFollowRoute(server)

	// Post Route
	v1.SetupPostRoute(server)

	// Comment Route
	v1.SetupCommentRoute(server)

	// Reply Route
	v1.SetupReplyRoute(server)

	// Like Route
	v1.SetupLikeRoute(server)

	// Hashtag Route
	v1.SetupHashtagRoute(server)

	// Mention Route
	v1.SetupMentionRoute(server)

	server.Run(":" + os.Getenv("PORT"))
}

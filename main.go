package main

import (
	"fmt"
	"historm_api/controllers"
	"historm_api/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	gin.SetMode(gin.ReleaseMode)

	var app *gin.Engine = gin.Default()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	app.Static("/files", "./public")

	app.GET("/api/v1/image/:file", func(ctx *gin.Context) {
		ctx.File("./public/display_images/" + ctx.Param("file"))
	})

	app.Use(middlewares.ApiMiddleware())
	app.POST("/api/v1/anime", controllers.RetrieveAnimeList)
	app.POST("/api/v1/latest-anime", controllers.RetrieveLatestAnimeList)
	app.POST("/api/v1/anime-detail", controllers.RetrieveAnimeDetail)
	app.POST("/api/v1/video", controllers.RetrieveVideo)

	user := app.Group("/api/v1/user")
	user.Use(middlewares.AuthMiddleware())
	{
		user.POST("/current", controllers.GetLoggedInUser)
		user.POST("/email", controllers.GetUserByEmail)
		user.POST("/create", controllers.InsertUser)
		user.POST("/delete", controllers.DeleteUser)
		user.POST("/update", controllers.UpdateUser)
		user.POST("/upload", controllers.UploadFile)
	}

	review := app.Group("/api/v1/review")
	review.Use(middlewares.AuthMiddleware())
	{
		review.POST("/create", controllers.CreateReview)
		review.POST("/likes/handle", controllers.HandleLikes)
		review.POST("/dislikes/handle", controllers.HandleDisLikes)
		review.POST("/retrieve/anime", controllers.RetrieveReviewsByAnime)
		review.POST("/retrieve/user", controllers.RetrieveReviewsByUser)
	}

	reply := app.Group("/api/v1/reply")
	reply.Use(middlewares.AuthMiddleware())
	{
		reply.POST("/create", controllers.CreateReply)
		reply.POST("/retrieve", controllers.RetrieveReplies)
		reply.POST("/retrieve/length", controllers.RetrieveRepliesLength)
		reply.POST("/likes/handle", controllers.HandleReplyLikes)
		reply.POST("/dislikes/handle", controllers.HandleReplyDisLikes)
	}

	favorite := app.Group("/api/v1/favorite")
	favorite.Use(middlewares.AuthMiddleware())
	{
		favorite.POST("/create", controllers.CreateFavorite)
		favorite.POST("/update", controllers.UpdateFavorite)
		favorite.POST("/retrieve", controllers.RetrieveFavorite)
		favorite.POST("/retrieve/all", controllers.RetrieveFavoriteAll)
	}

	rating := app.Group("/api/v1/rating")
	rating.Use(middlewares.AuthMiddleware())
	{
		rating.POST("/create", controllers.CreateRating)
		rating.POST("/update", controllers.UpdateRating)
		rating.POST("/retrieve/user", controllers.RetreiveRatingForUser)
		rating.POST("/retrieve/total", controllers.RetreiveTotalRating)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	fmt.Printf("Server is running in port %s", port)
	app.Run(":" + port)

}

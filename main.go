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

	app.Use(middlewares.ApiMiddleware())

	app.POST("/api/v1/anime", controllers.RetrieveAnimeList)
	app.POST("/api/v1/latest-anime", controllers.RetrieveLatestAnimeList)
	app.POST("/api/v1/anime-detail", controllers.RetrieveAnimeDetail)
	app.POST("/api/v1/video", controllers.RetrieveVideo)

	user := app.Group("/api/v1/user")
	user.Use(middlewares.AuthMiddleware())
	{
		user.POST("/", controllers.GetLoggedInUser)
		user.POST("/email", controllers.GetUserByEmail)
		user.POST("/create", controllers.InsertUser)
		user.POST("/delete", controllers.DeleteUser)
		user.POST("/update", controllers.UpdateUser)
	}

	review := app.Group("/api/v1/review")
	review.Use(middlewares.AuthMiddleware())
	{
		review.POST("/create", controllers.CreateReview)
		review.POST("/retrieve/anime", controllers.RetrieveReviewsByAnime)
		review.POST("/retrieve/user", controllers.RetrieveReviewsByUser)
	}

	reply := app.Group("/api/v1/reply")
	reply.Use(middlewares.AuthMiddleware())
	{
		reply.POST("/create", controllers.CreateReply)
		reply.POST("/retrieve", controllers.RetrieveReplies)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	fmt.Printf("Server is running in address http://%s:%s", host, port)
	app.Run(host + ":" + port)

}

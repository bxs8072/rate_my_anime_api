package middlewares

import (
	"context"
	"historm_api/configs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Token struct {
	Token string `bson:"token" json:"token"`
}

func AuthMiddleware() gin.HandlerFunc {
	var token Token

	return func(c *gin.Context) {
		if err := c.BindHeader(&token); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		app, _ := configs.FirebaseApp()

		client, err := app.Auth(context.Background())
		if err != nil {
			log.Fatalf("error verifying ID token: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		authToken, err := client.VerifyIDToken(context.Background(), token.Token)

		if err != nil {
			log.Fatalf("error verifying ID token: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("payload", bson.M{"uid": authToken.UID, "email": authToken.Claims["email"]})

		c.Next()
	}
}

package controllers

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadFile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userId := c.Request.FormValue("userId")

	file, err := c.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	extenstion := filepath.Ext(file.Filename)

	newFileName := userId + extenstion

	url := "public/display_images/" + newFileName

	if err := c.SaveUploadedFile(file, url); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": user}, bson.M{"$set": bson.M{"displayImage": "http://10.0.2.2:8080/api/v1/image/" + newFileName}})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"result": result,
		"data":   "https://rate-my-anime-api-git-oeghfd4pma-uc.a.run.app/api/v1/image/" + newFileName,
	})
}

package controllers

import (
	"context"
	"fmt"
	"historm_api/databases"
	"historm_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var listCollection *mongo.Collection = databases.GetCollection(databases.DB, "favorite")

func RetrieveFavorite(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := primitive.ObjectIDFromHex(body["user"])

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"user": userId, "anime.link": body["link"]}

	var favorites []models.Favorite

	cursor, err := listCollection.Find(ctx, filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)

	cursor.All(context.TODO(), &favorites)

	c.JSON(http.StatusAccepted, gin.H{"data": favorites})
}

func RetrieveFavoriteAll(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := primitive.ObjectIDFromHex(body["user"])

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"user": userId, "status": body["status"]}

	var favorites []models.Favorite

	cursor, err := listCollection.Find(ctx, filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)

	cursor.All(context.TODO(), &favorites)

	c.JSON(http.StatusAccepted, gin.H{"data": favorites})
}

func CreateFavorite(c *gin.Context) {
	var list models.Favorite

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0}

	result, err := listCollection.InsertOne(ctx, list)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

func UpdateFavorite(c *gin.Context) {

	var favorite models.Favorite

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	favorite.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0}

	filter := bson.M{"_id": bson.M{"$eq": favorite.Id}}
	update := bson.M{"$set": favorite}

	result, err := listCollection.UpdateOne(
		ctx, filter, update)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(result)
	c.JSON(http.StatusAccepted, result)
}

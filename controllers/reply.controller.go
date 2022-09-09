package controllers

import (
	"context"
	"historm_api/databases"
	"historm_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var replyCollection *mongo.Collection = databases.GetCollection(databases.DB, "reply")

func RetrieveReplies(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objID, err := primitive.ObjectIDFromHex(body["reviewId"])

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"reviewId": objID}
	var replies []models.Reply

	cursor, err := replyCollection.Find(ctx, filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cursor.All(context.TODO(), &replies)

	c.JSON(http.StatusAccepted, replies)
}

func CreateReply(c *gin.Context) {
	var reply models.Reply

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&reply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0}

	result, err := replyCollection.InsertOne(ctx, reply)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

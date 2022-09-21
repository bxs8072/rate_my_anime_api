package controllers

import (
	"context"
	"historm_api/databases"
	"historm_api/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ratingCollection *mongo.Collection = databases.GetCollection(databases.DB, "rating")

func trimLink(link string) string {
	var last string = strings.Split(strings.Split(link, "gogoanime.")[1], "/")[0]

	return last
}

func RetreiveRatingForUser(c *gin.Context) {
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

	filter := bson.M{"user": userId, "anime.link": body["animeLink"]}

	var ratings []models.Rating

	cursor, err := ratingCollection.Find(ctx, filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)

	cursor.All(context.TODO(), &ratings)

	c.JSON(http.StatusAccepted, gin.H{"data": ratings})
}

func RetreiveTotalRating(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trimLink(body["animeLink"])

	matchStage := bson.D{{Key: "$match", Value: bson.M{"anime.link": body["animeLink"]}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": "$anime", "total": bson.M{"$avg": "$rating"}}}}

	showInfoCursor, err := ratingCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		panic(err)
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		panic(err)
	}

	c.JSON(http.StatusAccepted, gin.H{"data": showsWithInfo})
}

func CreateRating(c *gin.Context) {
	var rating models.Rating

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0}

	result, err := ratingCollection.InsertOne(ctx, rating)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

func UpdateRating(c *gin.Context) {
	var rating models.Rating

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0}

	filter := bson.M{"_id": bson.M{"$eq": rating.Id}}
	update := bson.M{"$set": rating}

	result, err := ratingCollection.UpdateOne(
		ctx, filter, update)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

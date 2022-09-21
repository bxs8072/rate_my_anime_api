package controllers

import (
	"context"
	"historm_api/databases"
	"historm_api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var reviewCollection *mongo.Collection = databases.GetCollection(databases.DB, "review")

func RetrieveReviewsByAnime(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reviewList []map[string]interface{}
	var populate bson.M = bson.M{"$lookup": bson.M{
		"from":         "user",
		"localField":   "user",
		"foreignField": "_id",
		"as":           "user",
	}}

	var search bson.M = bson.M{"$match": bson.M{"anime.link": body["animeLink"]}}
	var sort bson.M = bson.M{"$sort": bson.M{"createdAt": -1}}
	var limit bson.M = bson.M{"$limit": 2}
	skipValue, _ := strconv.Atoi(body["skipValue"])

	var skip bson.M = bson.M{"$skip": skipValue}

	cursor, err := reviewCollection.Aggregate(ctx, []bson.M{search, sort, skip, limit, populate})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cursor.All(context.TODO(), &reviewList)

	c.JSON(http.StatusAccepted, reviewList)
}

func RetrieveReviewsByUser(c *gin.Context) {
	var body map[string]string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objID, err := primitive.ObjectIDFromHex(body["userId"])

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"userId": objID}

	var reviews []models.Review
	cursor, err := reviewCollection.Find(ctx, filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)

	cursor.All(context.TODO(), &reviews)

	c.JSON(http.StatusAccepted, reviews)
}

func CreateReview(c *gin.Context) {
	var review models.Review

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0}

	result, err := reviewCollection.InsertOne(ctx, review)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

func DeleteReview(c *gin.Context) {
	var body map[string]string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviewId, err := primitive.ObjectIDFromHex(body["reviewId"])

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := primitive.ObjectIDFromHex(body["userId"])

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := reviewCollection.DeleteOne(ctx, bson.M{"reviewId": reviewId, "userId": userId})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, result)
}

func HandleLikes(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := primitive.ObjectIDFromHex(body["user"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := primitive.ObjectIDFromHex(body["id"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ones []models.Review

	cursor, err := reviewCollection.Find(ctx, bson.M{"_id": review, "likes": bson.M{"$in": []primitive.ObjectID{user}}})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cursor.All(ctx, &ones)

	var result *mongo.UpdateResult

	if ones == nil {
		result, err = reviewCollection.UpdateByID(ctx, review, bson.M{"$push": bson.M{"likes": user}})
	} else {
		result, err = reviewCollection.UpdateByID(ctx, review, bson.M{"$pull": bson.M{"likes": user}})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

func HandleDisLikes(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := primitive.ObjectIDFromHex(body["user"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := primitive.ObjectIDFromHex(body["id"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ones []models.Review

	cursor, err := reviewCollection.Find(ctx, bson.M{"_id": review, "dislikes": bson.M{"$in": []primitive.ObjectID{user}}})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cursor.All(ctx, &ones)

	var result *mongo.UpdateResult

	if ones == nil {
		result, err = reviewCollection.UpdateByID(ctx, review, bson.M{"$push": bson.M{"dislikes": user}})
	} else {
		result, err = reviewCollection.UpdateByID(ctx, review, bson.M{"$pull": bson.M{"dislikes": user}})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

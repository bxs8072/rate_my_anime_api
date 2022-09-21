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

var replyCollection *mongo.Collection = databases.GetCollection(databases.DB, "reply")

func RetrieveReplies(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var replyList []map[string]interface{}

	var userPopulate bson.M = bson.M{"$lookup": bson.M{
		"from":         "user",
		"localField":   "user",
		"foreignField": "_id",
		"as":           "user",
	}}

	review, err := primitive.ObjectIDFromHex(body["reviewId"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var search bson.M = bson.M{"$match": bson.M{"review": review}}
	var sort bson.M = bson.M{"$sort": bson.M{"createdAt": -1}}
	var limit bson.M = bson.M{"$limit": 5}
	skipValue, _ := strconv.Atoi(body["skipValue"])

	var skip bson.M = bson.M{"$skip": skipValue}

	cursor, err := replyCollection.Aggregate(ctx, []bson.M{search, sort, skip, limit, userPopulate})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cursor.All(context.TODO(), &replyList)

	c.JSON(http.StatusAccepted, replyList)
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

func RetrieveRepliesLength(c *gin.Context) {
	var body map[string]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var replyList []map[string]interface{}

	var userPopulate bson.M = bson.M{"$lookup": bson.M{
		"from":         "user",
		"localField":   "user",
		"foreignField": "_id",
		"as":           "user",
	}}

	review, err := primitive.ObjectIDFromHex(body["reviewId"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var search bson.M = bson.M{"$match": bson.M{"review": review}}
	var limit bson.M = bson.M{"$limit": 10}

	cursor, err := replyCollection.Aggregate(ctx, []bson.M{search, limit, userPopulate})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cursor.All(context.TODO(), &replyList)

	c.JSON(http.StatusAccepted, bson.M{"length": len(replyList)})
}

func HandleReplyLikes(c *gin.Context) {
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

	reply, err := primitive.ObjectIDFromHex(body["id"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ones []models.Review

	cursor, err := replyCollection.Find(ctx, bson.M{"_id": reply, "likes": bson.M{"$in": []primitive.ObjectID{user}}})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cursor.All(ctx, &ones)

	var result *mongo.UpdateResult

	if ones == nil {
		result, err = replyCollection.UpdateByID(ctx, reply, bson.M{"$push": bson.M{"likes": user}})
	} else {
		result, err = replyCollection.UpdateByID(ctx, reply, bson.M{"$pull": bson.M{"likes": user}})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

func HandleReplyDisLikes(c *gin.Context) {
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

	reply, err := primitive.ObjectIDFromHex(body["id"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ones []models.Review

	cursor, err := replyCollection.Find(ctx, bson.M{"_id": reply, "dislikes": bson.M{"$in": []primitive.ObjectID{user}}})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cursor.All(ctx, &ones)

	var result *mongo.UpdateResult

	if ones == nil {
		result, err = replyCollection.UpdateByID(ctx, reply, bson.M{"$push": bson.M{"dislikes": user}})
	} else {
		result, err = replyCollection.UpdateByID(ctx, reply, bson.M{"$pull": bson.M{"dislikes": user}})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

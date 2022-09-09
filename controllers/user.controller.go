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

var userCollection *mongo.Collection = databases.GetCollection(databases.DB, "user")

// POST => /api/v1/user/retrieve
func GetLoggedInUser(c *gin.Context) {
	payload, exist := c.Get("payload")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization"})
		return
	}

	var uid = payload.(primitive.M)["uid"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"uid": uid}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, user)
}

// POST => /api/v1/user/retrieve/email
func GetUserByEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var body map[string]string
	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	err := userCollection.FindOne(ctx, bson.M{"email": body["email"]}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// POST => /api/v1/user/retrieve/create
func InsertUser(c *gin.Context) {
	payload, exist := c.Get("payload")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization"})
		return
	}

	var uid = payload.(primitive.M)["uid"]
	var email = payload.(primitive.M)["email"]

	fmt.Println(email)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var body map[string]string
	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var firstName = body["firstName"]
	var middleName = body["middleName"]
	var lastName = body["lastName"]
	var displayImage = body["displayImage"]

	fmt.Println(firstName)

	result, err := userCollection.InsertOne(ctx, bson.M{
		"uid":          uid,
		"email":        email,
		"firstName":    firstName,
		"middleName":   middleName,
		"lastName":     lastName,
		"displayImage": displayImage,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": result.InsertedID,
	})
}

// POST => /api/v1/user/update
func UpdateUser(c *gin.Context) {
	payload, exist := c.Get("payload")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization"})
		return
	}

	var uid = payload.(primitive.M)["uid"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := userCollection.UpdateOne(ctx, bson.D{{Key: "uid", Value: uid}}, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": result.UpsertedID,
	})
}

// POST => /api/v1/user/delete
func DeleteUser(c *gin.Context) {
	payload, exist := c.Get("payload")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization"})
		return
	}

	var uid = payload.(primitive.M)["uid"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.D{{Key: "uid", Value: uid}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"data": result.DeletedCount,
	})
}

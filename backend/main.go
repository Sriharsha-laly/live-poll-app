package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Poll struct {
	ID       string   `json:"id" bson:"_id,omitempty"`
	Question string   `json:"question" bson:"question"`
	Options  []string `json:"options" bson:"options"`
}

var collection *mongo.Collection

func main() {

	// Get MongoDB connection string from environment variable
	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	// Test MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	log.Println("MongoDB Connected Successfully!")

	// Database and collection
	collection = client.Database("live_poll_db").Collection("polls")

	// Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// Home route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Backend Running",
		})
	})

	// Create poll
	router.POST("/create-poll", func(c *gin.Context) {

		var poll Poll

		if err := c.ShouldBindJSON(&poll); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Data",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := collection.InsertOne(ctx, poll)

		if err != nil {
			log.Println("Insert error:", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save poll",
			})
			return
		}

		// Convert MongoDB ObjectID to string
		if objectID, ok := result.InsertedID.(bson.ObjectID); ok {
			poll.ID = objectID.Hex()
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Poll Created Successfully",
			"poll":    poll,
		})
	})

	// Get all polls
	router.GET("/polls", func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch polls",
			})
			return
		}

		defer cursor.Close(ctx)

		var polls []Poll

		if err := cursor.All(ctx, &polls); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to decode polls",
			})
			return
		}

		if polls == nil {
			polls = []Poll{}
		}

		c.JSON(http.StatusOK, polls)
	})

	// Render provides PORT through an environment variable
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
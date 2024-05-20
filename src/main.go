package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	mongoClient *mongo.Client
	redisClient *redis.Client
)

const (
	mongoURI       = "mongodb://localhost:27017/?connectTimeoutMS=10000"
	redisAddr      = "localhost:6379"
	imageUploadDir = "./uploads/"
	bcryptCost     = 6
)

type User struct {
	Username     string `json:"username" bson:"username" validate:"required"`
	Password     string `json:"password" bson:"-"`
	PasswordHash string `json:"-" bson:"password_hash"`
	DisplayName  string `json:"display_name" bson:"display_name"`
	Relation     string `json:"relation" bson:"relation"`
}

type DiaryEntry struct {
	ID         string    `json:"id" bson:"id"`
	Username   string    `json:"username" bson:"username"`
	Content    string    `json:"content" bson:"content"`
	ImagePaths []string  `json:"image_paths,omitempty" bson:"image_paths,omitempty"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
}

func main() {
	// Initialize MongoDB
	var err error
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.TODO())

	// Initialize Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer redisClient.Close()

	// Test Redis Connection
	_, err = redisClient.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// Create Gin router
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.POST("/register", registerHandler)
	r.POST("/diary", authMiddleware, createDiaryEntry)
	r.GET("/diary", authMiddleware, getDiaryEntries)

	// Serve uploaded images
	r.Static("/uploads", imageUploadDir)

	err = r.Run(":8818")
	if err != nil {
		panic(err)
	}
}

func loginHandler(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	collection := mongoClient.Database("baby_diary").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"username": credentials.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username1 or password"})
		return
	}
	fmt.Println("user:", user, "credential:", credentials)

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		fmt.Println("err:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password1"})
		return
	}

	token, err := generateToken(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
func registerHandler(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	collection := mongoClient.Database("baby_diary").Collection("users")

	// Check if username already exists
	var existingUser User
	err := collection.FindOne(context.TODO(), bson.M{"username": newUser.Username}).Decode(&existingUser)
	if err == nil {
		// User with this username already exists
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcryptCost)
	fmt.Println("Checking password: ", passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	newUser.PasswordHash = string(passwordHash)

	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func createDiaryEntry(c *gin.Context) {
	username := c.MustGet("username").(string)

	// 判断 Content-Type 是否是 application/json
	if c.ContentType() == "application/json" {
		var entry struct {
			Content string `json:"content"`
		}
		if err := c.BindJSON(&entry); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		diaryEntry := DiaryEntry{
			ID:        uuid.New().String(),
			Username:  username,
			Content:   entry.Content,
			Timestamp: time.Now(),
		}

		collection := mongoClient.Database("baby_diary").Collection("diary_entries")
		_, err := collection.InsertOne(context.TODO(), diaryEntry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create diary entry"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Diary entry created successfully"})
	} else {
		// 处理 multipart/form-data 请求
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form"})
			return
		}

		files := form.File["images"]
		if len(files) > 9 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Too many images, maximum is 9"})
			return
		}

		var entry struct {
			Content string `json:"content"`
		}
		entry.Content = form.Value["content"][0]

		imagePaths := []string{}
		for _, file := range files {
			imageID := uuid.New().String()
			imagePath := filepath.Join(imageUploadDir, imageID+"-"+file.Filename)
			if err := c.SaveUploadedFile(file, imagePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save image"})
				return
			}
			imagePaths = append(imagePaths, imagePath)
		}

		diaryEntry := DiaryEntry{
			ID:         uuid.New().String(),
			Username:   username,
			Content:    entry.Content,
			ImagePaths: imagePaths,
			Timestamp:  time.Now(),
		}

		collection := mongoClient.Database("baby_diary").Collection("diary_entries")
		_, err = collection.InsertOne(context.TODO(), diaryEntry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create diary entry"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Diary entry created successfully"})
	}
}

func getDiaryEntries(c *gin.Context) {
	collection := mongoClient.Database("baby_diary").Collection("diary_entries")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch diary entries"})
		return
	}
	defer cursor.Close(context.TODO())

	var entries []DiaryEntry
	if err := cursor.All(context.TODO(), &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not decode diary entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		c.Abort()
		return
	}

	// 去掉 "Bearer " 前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.Abort()
		return
	}

	username, err := validateToken(token)
	if err != nil {
		log.Printf("Token validation error: %v, username: %s", err, username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	log.Printf("Authenticated user %s with token %s", username, token)
	c.Set("username", username)
	c.Next()
}

func generateToken(username string) (string, error) {
	token := uuid.New().String()
	err := redisClient.Set(context.TODO(), token, username, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to store token in Redis: %v", err)
		return "", err
	}
	log.Printf("Stored token for user %s: %s", username, token)
	return token, nil
}

func validateToken(token string) (string, error) {
	username, err := redisClient.Get(context.TODO(), token).Result()
	if errors.Is(err, redis.Nil) {
		log.Printf("Token not found in Redis: %s", token)
		return "", fmt.Errorf("token not found")
	} else if err != nil {
		log.Printf("Error retrieving token from Redis: %v", err)
		return "", err
	}
	return username, nil
}

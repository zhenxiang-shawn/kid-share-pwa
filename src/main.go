// package main
//
// import (
//
//	"context"
//	"fmt"
//	"go.mongodb.org/mongo-driver/mongo/readpref"
//	"kid-share/controller"
//	"kid-share/model"
//
//	"github.com/gin-gonic/gin"
//
// )
//
//	func main() {
//		ctx := context.Background()
//		r := gin.Default()
//
//		// 连接到 MongoDB
//		if err := model.Connect(ctx, "mongodb://localhost:27017/?connectTimeoutMS=10000"); err != nil {
//			panic(err)
//		}
//		// 测试链接
//		// 使用客户端实例进行一个操作，如 Ping，来测试连接
//		err := model.Client.Ping(ctx, readpref.Primary())
//		if err != nil {
//			fmt.Println("Failed to ping MongoDB:", err)
//			panic(err)
//		}
//
//		v1 := r.Group("/api/v1")
//		{
//			userCtrl := controller.UserController{}
//			v1.POST("/user", userCtrl.Create)
//			v1.GET("/user/:username", userCtrl.GetByUsername)
//			v1.PUT("/user/:username", userCtrl.UpdateByUsername)
//			v1.DELETE("/user", userCtrl.DeleteByID)
//		}
//
//		r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
//	}
package main

import (
	"context"
	"log"
	"net/http"
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
)

type User struct {
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password_hash"`
	DisplayName  string `json:"display_name" bson:"display_name"`
	Relation     string `json:"relation" bson:"relation"`
}

type DiaryEntry struct {
	ID        string    `json:"id" bson:"id"`
	Username  string    `json:"username" bson:"username"`
	Content   string    `json:"content" bson:"content"`
	ImagePath string    `json:"image_path,omitempty" bson:"image_path,omitempty"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
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
	//redisClient = redis.NewClient(&redis.Options{
	//	Addr: redisAddr,
	//})
	//defer redisClient.Close()

	// Create Gin router
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.POST("/register", registerHandler)
	r.POST("/diary", authMiddleware, createDiaryEntry)
	r.GET("/diary", authMiddleware, getDiaryEntries)

	// Serve uploaded images
	r.Static("/uploads", imageUploadDir)

	r.Run(":8818")
}

func loginHandler(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
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

	collection := mongoClient.Database("baby_diary").Collection("users")

	// Check if username already exists
	var existingUser User
	err := collection.FindOne(context.TODO(), bson.M{"username": newUser.Username}).Decode(&existingUser)
	if err == nil {
		// User with this username already exists
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
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
/*
func registerHandler(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	newUser.PasswordHash = string(passwordHash)

	collection := mongoClient.Database("baby_diary").Collection("users")
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
*/
func createDiaryEntry(c *gin.Context) {
	username := c.MustGet("username").(string)
	var entry struct {
		Content string `json:"content"`
	}
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	file, err := c.FormFile("image")
	imagePath := ""
	if err == nil {
		imageID := uuid.New().String()
		imagePath = imageUploadDir + imageID + "-" + file.Filename
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save image"})
			return
		}
	}

	diaryEntry := DiaryEntry{
		ID:        uuid.New().String(),
		Username:  username,
		Content:   entry.Content,
		ImagePath: imagePath,
		Timestamp: time.Now(),
	}

	collection := mongoClient.Database("baby_diary").Collection("diary_entries")
	_, err = collection.InsertOne(context.TODO(), diaryEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create diary entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diary entry created successfully"})
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

	username, err := validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("username", username)
	c.Next()
}

func generateToken(username string) (string, error) {
	token := uuid.New().String()
	err := redisClient.Set(context.TODO(), token, username, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func validateToken(token string) (string, error) {
	username, err := redisClient.Get(context.TODO(), token).Result()
	if err != nil {
		return "", err
	}
	return username, nil
}

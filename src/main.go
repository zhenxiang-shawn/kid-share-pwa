package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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
	Avatar       string `json:"avatar" bson:"avatar"`
	Relation     string `json:"relation" bson:"relation"`
}

type DiaryEntry struct {
	ID         string    `json:"id" bson:"id"`
	Username   string    `json:"username" bson:"username"`
	Content    string    `json:"content" bson:"content"`
	ImagePaths []string  `json:"image_paths,omitempty" bson:"image_paths,omitempty"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
}

type DiaryResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
	Data    struct {
		Diaries []DiaryEntry `json:"diaries"`
		Total   int          `json:"total"`
	} `json:"data"`
}

type LoginResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
	Data    struct {
		Token       string `json:"token"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Avatar      string `json:"avatar"`
		Relation    string `json:"relation"`
	} `json:"data"`
}

type UserInfoResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
	Data    struct {
		Name     string `json:"name"`
		Avatar   string `json:"avatar"`
		Relation string `json:"relation"`
	} `json:"data"`
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
	apiGroup := r.Group("api")
	{
		apiGroup.POST("/user/login", loginHandler)
		apiGroup.GET("/user/:username", getUserInfo)
		apiGroup.POST("/user/register", registerHandler)
		apiGroup.POST("/diary", authMiddleware, createDiaryEntry)
		apiGroup.GET("/diaries", authMiddleware, getDiaryEntries)
	}

	// Serve uploaded images
	apiGroup.Static("/uploads", imageUploadDir)

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
	response := LoginResponseData{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Login successful"),
		Ok:      true,
	}
	response.Data.Username = user.Username
	response.Data.Token = token
	response.Data.DisplayName = user.DisplayName
	c.JSON(http.StatusOK, response)

}

func getUserInfo(c *gin.Context) {
	//username := c.MustGet("username").(string)
	username := c.Param("username")

	collection := mongoClient.Database("baby_diary").Collection("users")
	var user = User{}
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "ok": false})
		return
	}
	var userInfo = UserInfoResponseData{}
	userInfo.Ok = true
	userInfo.Data.Avatar = user.Avatar
	userInfo.Data.Name = user.DisplayName
	userInfo.Data.Relation = user.Relation

	c.JSON(http.StatusOK, userInfo)

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
			Content    string   `json:"content"`
			ImagePaths []string `json:"image_paths"`
		}
		if err := c.BindJSON(&entry); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		diaryEntry := DiaryEntry{
			ID:         uuid.New().String(),
			Username:   username,
			Content:    entry.Content,
			ImagePaths: entry.ImagePaths,
			Timestamp:  time.Now(),
		}

		collection := mongoClient.Database("baby_diary").Collection("diary_entries")
		_, err := collection.InsertOne(context.TODO(), diaryEntry)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Could not create diary entry",
					"ok":      false,
				})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Diary entry created successfully",
			"ok":      true,
		})
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

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Diary entry created successfully",
			"ok":      true,
		})
	}
	//} else {
	//	// 处理 multipart/form-data 请求
	//	form, err := c.MultipartForm()
	//	if err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form"})
	//		return
	//	}
	//
	//	files := form.File["images"]
	//	if len(files) > 9 {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many images, maximum is 9"})
	//		return
	//	}
	//
	//	var entry struct {
	//		Content string `json:"content"`
	//	}
	//	entry.Content = form.Value["content"][0]
	//
	//	imagePaths := []string{}
	//	for _, file := range files {
	//		imageID := uuid.New().String()
	//		imagePath := filepath.Join(imageUploadDir, imageID+"-"+file.Filename)
	//		if err := c.SaveUploadedFile(file, imagePath); err != nil {
	//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save image"})
	//			return
	//		}
	//		imagePaths = append(imagePaths, imagePath)
	//	}
	//
	//	diaryEntry := DiaryEntry{
	//		ID:         uuid.New().String(),
	//		Username:   username,
	//		Content:    entry.Content,
	//		ImagePaths: imagePaths,
	//		Timestamp:  time.Now(),
	//	}
	//
	//	collection := mongoClient.Database("baby_diary").Collection("diary_entries")
	//	_, err = collection.InsertOne(context.TODO(), diaryEntry)
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create diary entry"})
	//		return
	//	}
	//
	//	c.JSON(http.StatusOK, gin.H{"message": "Diary entry created successfully"})
	//}
}

func getDiaryEntries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	skip := (page - 1) * limit

	collection := mongoClient.Database("baby_diary").Collection("diary_entries")
	filter := bson.M{}

	var diaries []DiaryEntry
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{"date", -1}}) // Sort by date in descending order

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch diary entries"})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var diary DiaryEntry
		if err := cursor.Decode(&diary); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode diary entry"})
			return
		}
		diaries = append(diaries, diary)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}
	total, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "faied to count documents"})
		return
	}
	response := DiaryResponseData{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Diary entries retrieved successfully"),
		Ok:      true,
	}
	response.Data.Diaries = diaries
	response.Data.Total = int(total)

	c.JSON(http.StatusOK, response)
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

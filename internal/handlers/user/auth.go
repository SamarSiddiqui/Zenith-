package user

import (
	"context"
	"net/http"
	"time"

	"habit-tracker/internal/db"
	"habit-tracker/internal/models"
	"habit-tracker/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Signup handles new user registration
func Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	collection := db.GetCollection("habit-tracker", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email already in use", nil)
		return
	} else if err != mongo.ErrNoDocuments {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to encrypt password", err)
		return
	}

	// Create arbitrary ObjectID
	userID := primitive.NewObjectID()

	token, err := utils.GenerateToken(userID.Hex())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	// Create user
	newUser := models.User{
		ID:          userID,
		Email:       req.Email,
		Password:    hashedPassword,
		AccessToken: token, // Store the token in the DB
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// Return User and Token
	res := models.AuthResponse{
		Token: token,
		User:  newUser,
	}

	c.JSON(http.StatusCreated, res)
}

// Login handles user authentication
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	collection := db.GetCollection("habit-tracker", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	// Generate new JWT Token upon login
	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	// Update the access_token in database
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"access_token": token,
			"updated_at":   time.Now(),
		}},
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update access token in db", err)
		return
	}

	// Make sure the object returned reflects the new token and time
	user.AccessToken = token
	user.UpdatedAt = time.Now()

	res := models.AuthResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusOK, res)
}

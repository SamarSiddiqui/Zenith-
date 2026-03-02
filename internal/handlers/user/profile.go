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

// GetProfile handles retrieving the user's profile
func GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	collection := db.GetCollection("habit-tracker", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}

	// Initialize empty structs if nil to match expected response shape
	profile := user.Profile
	if profile == nil {
		profile = &models.UserProfile{}
	}
	preferences := user.Preferences
	if preferences == nil {
		preferences = &models.UserPreferences{}
	}
	actions := user.Actions
	if actions == nil {
		actions = &models.UserActions{}
	}

	c.JSON(http.StatusOK, gin.H{
		"profile":     profile,
		"preferences": preferences,
		"actions":     actions,
	})
}

// PatchProfileRequest defines the payload for patching the user profile
type PatchProfileRequest struct {
	Profile     map[string]interface{} `json:"profile,omitempty"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
	Actions     map[string]interface{} `json:"actions,omitempty"`
}

// PatchProfile handles updating the user's profile, preferences, and actions
func PatchProfile(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req PatchProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	collection := db.GetCollection("habit-tracker", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateFields := bson.M{}

	if len(req.Profile) > 0 {
		for k, v := range req.Profile {
			updateFields["profile."+k] = v
		}
	}
	if len(req.Preferences) > 0 {
		for k, v := range req.Preferences {
			updateFields["preferences."+k] = v
		}
	}
	if len(req.Actions) > 0 {
		for k, v := range req.Actions {
			updateFields["actions."+k] = v
		}
	}

	if len(updateFields) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "No valid fields provided for update", nil)
		return
	}

	updateFields["updated_at"] = time.Now()

	update := bson.M{
		"$set": updateFields,
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user profile", err)
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	// Fetch updated user to respond with new state
	var updatedUser models.User
	err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&updatedUser)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch updated user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Profile updated successfully",
		"profile":     updatedUser.Profile,
		"preferences": updatedUser.Preferences,
		"actions":     updatedUser.Actions,
	})
}

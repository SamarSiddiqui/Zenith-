package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the User model in the database
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"-"` // Don't return password in JSON
	AccessToken string             `bson:"access_token" json:"access_token,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// SignupRequest defines the payload for creating a new user
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest defines the payload for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse defines the payload returned upon successful auth
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

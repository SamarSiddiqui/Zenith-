package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserProfile represents the nested profile data
type UserProfile struct {
	DisplayName string `bson:"display_name,omitempty" json:"display_name,omitempty"`
	Username    string `bson:"username,omitempty" json:"username,omitempty"`
	Email       string `bson:"email,omitempty" json:"email,omitempty"`
}

// UserPreferences represents the user's preferences
type UserPreferences struct {
	StartOfWeek  string `bson:"startOfWeek,omitempty" json:"startOfWeek,omitempty"`
	ReminderTime string `bson:"reminderTime,omitempty" json:"reminderTime,omitempty"`
	VacationMode bool   `bson:"vacationMode,omitempty" json:"vacationMode,omitempty"`
	Theme        string `bson:"theme,omitempty" json:"theme,omitempty"`
}

// UserActions represents specific actions/flags for the user
type UserActions struct {
	DeleteAccount bool `bson:"deleteAccount,omitempty" json:"deleteAccount,omitempty"`
}

// User represents the User model in the database
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"-"` // Don't return password in JSON
	AccessToken string             `bson:"access_token,omitempty" json:"access_token,omitempty"`
	Profile     *UserProfile       `bson:"profile,omitempty" json:"profile,omitempty"`
	Preferences *UserPreferences   `bson:"preferences,omitempty" json:"preferences,omitempty"`
	Actions     *UserActions       `bson:"actions,omitempty" json:"actions,omitempty"`
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

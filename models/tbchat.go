package models

import "time"

type Chat struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	ChatroomID string    `json:"chatroom_id"`
	Message    string    `json:"message"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	UserID     string    `json:"user_id"`
	CreateBy   string    `json:"create_by"`
	CreateDate time.Time `json:"create_date"`
	IsRead     bool      `json:"is_read"`
}

package models

import "time"

type SocialStoryUser struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint64    `json:"user_id"`
	StoryID   uint64    `json:"story_id"`
	User      User      `json:"user" grom:"references:UserID"`
}

type SocialStoryUserSeen struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	UserID  uint64 `json:"user_id"`
	StoryID uint64 `json:"story_id"`
}

func (SocialStoryUserSeen) TableName() string {
	return "social_story_users"
}

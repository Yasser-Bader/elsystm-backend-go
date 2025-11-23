package models

import "time"

type SocialStoryMigration struct {
	ID                 uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	UserId             uint64    `json:"user_id"`
	Content            string    `json:"content"`
	Duration           string    `json:"duration"`
	FileName           string    `json:"file_name"`
	FileNameAdditional string    `json:"file_name_additional"`
	FileType           string    `json:"type"`
}

type SocialStoryCU struct {
	ID                 uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	UserId             uint64    `json:"user_id"`
	Content            string    `json:"content"`
	Duration           string    `json:"duration"`
	FileName           string    `json:"file_name"`
	FileNameAdditional string    `json:"file_name_additional"`
	FileType           string    `json:"type"`
}
type SocialStory struct {
	ID                 uint64    `json:"id" gorm:"primaryKey"`
	UserId             uint64    `json:"user_id"`
	Content            string    `json:"content"`
	Duration           string    `json:"duration"`
	FileName           string    `json:"file_name"`
	FileNameAdditional string    `json:"file_name_additional"`
	FileType           string    `json:"type"`
	CountSeen          int64     `json:"count_seen"`
	CreatedAt          time.Time `json:"created_at"`
}

type SocialStoryWithoutSeenSerializer struct {
	ID                 uint64              `json:"id" gorm:"primaryKey"`
	UserId             uint64              `json:"user_id"`
	Content            string              `json:"content"`
	Duration           string              `json:"duration"`
	FileName           string              `json:"file_name"`
	FileNameAdditional string              `json:"file_name_additional"`
	FileType           string              `json:"type"`
	CreatedAt          time.Time           `json:"created_at"`
	CountSeen          int64               `json:"count_seen"`
	Seen               SocialStoryUserSeen `json:"seen" gorm:"foreignKey:StoryID;"`
}

func (SocialStoryMigration) TableName() string {
	return "social_stories"
}

func (SocialStoryWithoutSeenSerializer) TableName() string {
	return "social_stories"
}

func (SocialStoryCU) TableName() string {
	return "social_stories"
}

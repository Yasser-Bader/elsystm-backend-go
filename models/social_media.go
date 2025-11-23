package models

import "time"

type SocialMedia struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uint64    `json:"user_id"`
	ModelId   uint64    `json:"model_id"`
	ModelType string    `json:"model_type" gorm:"varchar(100)"`
	FileName  string    `json:"file_name"`
	FileType  string    `json:"type"`
}

package models

import "time"

type SocialHashtagPCR struct {
	ID        uint64        `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	HashtagId uint64        `json:"hashtag_id"`
	ModelId   uint64        `json:"model_id"`
	ModelType string        `json:"model_type" gorm:"varchar(100)"`
	Hashtag   SocialHashtag `json:"hashtag" gorm:"foreignKey:HashtagId"`
}

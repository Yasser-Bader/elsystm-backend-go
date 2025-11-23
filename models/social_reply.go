package models

import "time"

type SocialReplyMigration struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uint64    `json:"user_id"`
	PostId    uint64    `json:"post_id"`
	CommentId uint64    `json:"comment_id"`
	Content   string    `json:"content"`
}

type SocialReply struct {
	ID        uint64          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	UserId    uint64          `json:"user_id"`
	PostId    uint64          `json:"post_id"`
	CommentId uint64          `json:"comment_id"`
	Content   string          `json:"content"`
	Media     SocialMedia     `json:"media" gorm:"foreignKey:ModelId"`
	Hashtags  []SocialHashtag `json:"hashtags" gorm:"foreignKey:ModelId"`
	Mentions  []SocialMention `json:"mentions" gorm:"foreignKey:ModelId"`
	Likes     []SocialLike    `json:"likes" gorm:"foreignKey:ModelId"`
}

func (SocialReplyMigration) TableName() string {
	return "social_replies"
}

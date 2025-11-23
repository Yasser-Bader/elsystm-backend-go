package models

import "time"

type SocialPostMigration struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uint64    `json:"user_id"`
	Content   string    `json:"content"`
}

type SocialPostShare struct {
	ID          uint64          `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	PostId      uint64          `json:"post_id"`
	PostShareId uint64          `json:"post_share_id"`
	PostShared  SocialPostShard `json:"post" gorm:"foreignKey:PostShareId"`
}

type SocialPost struct {
	ID        uint64        `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	UserId    uint64        `json:"user_id"`
	Content   string        `json:"content"`
	Media     []SocialMedia `json:"media" gorm:"foreignKey:ModelId"`
}

type SocialPostShard struct {
	ID        uint64        `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	UserId    uint64        `json:"user_id"`
	Content   string        `json:"content"`
	Media     []SocialMedia `json:"media" gorm:"foreignKey:ModelId"`
	User      User          `json:"user" gorm:"foreignKey:UserId"`
}

type SocialPostHome struct {
	ID           uint64             `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserId       uint64             `json:"user_id"`
	Content      string             `json:"content"`
	Media        []SocialMedia      `json:"media" gorm:"foreignKey:ModelId"`
	User         User               `json:"user" gorm:"foreignKey:UserId"`
	SharedPost   SocialPostShare    `json:"shared_post" gorm:"foreignKey:PostId"`
	Mentions     []SocialMention    `json:"mentions" gorm:"foreignKey:ModelId"`
	Hashtags     []SocialHashtagPCR `json:"hashtags" gorm:"foreignKey:ModelId"`
	CountComment int64              `json:"count_comment"`
	CountLike    int64              `json:"count_like"`
	CountShare   int64              `json:"count_share"`
}

type SocialGetPosts struct {
	Data     interface{} `json:"data"`
	Comments int64       `json:"comments"`
	Likes    int64       `json:"likes"`
}

func (SocialPostMigration) TableName() string {
	return "social_posts"
}

func (SocialPostHome) TableName() string {
	return "social_posts"
}

func (SocialPostShard) TableName() string {
	return "social_posts"
}

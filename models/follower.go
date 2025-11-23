package models

import (
	"time"
)

type Follower struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserId      uint64    `json:"user_id"`
	FollowingId uint64    `json:"following_id"`
	Status      bool      `json:"status"`
	Following   User      `json:"following" grom:"references:FollowingId"`
}

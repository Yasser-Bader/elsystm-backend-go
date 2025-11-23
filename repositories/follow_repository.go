package repositories

import (
	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"gorm.io/gorm"
)

type FollowRepository interface {
	Followers(page int, userId uint64, pagination *util.Pagination)
	FollowingID(userId uint64) ([]uint64, error)
}

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository() FollowRepository {
	return &followRepository{db: config.ConnectDB()}
}

func (r followRepository) Followers(page int, userId uint64, pagination *util.Pagination) {
	var followers []models.Follower
	var count int64

	// Count rows
	r.db.Model(&followers).Where("user_id=?", userId).Where("status=?", 1).Where("following_id <> ?", userId).Count(&count)
	// Get Followers
	r.db.Scopes(util.Paginate(page, pagination, count, 10)).Preload("Following").Preload("Following.Avatar", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "App\\Domain\\Users\\Entity\\Models\\User").Where("collection_name=?", "avatar")
	}).Where("user_id=?", userId).Where("status=?", 1).Where("following_id <> ?", userId).Find(&followers)

	pagination.Data = followers
}

func (r followRepository) FollowingID(userId uint64) ([]uint64, error) {
	var followers []models.Follower
	var userIds []uint64
	if err := r.db.Model(&followers).Where("user_id=?", userId).Where("status=?", 1).Where("following_id <> ?", userId).Pluck("following_id", &userIds).Error; err != nil {
		return userIds, err
	}
	return userIds, nil
}

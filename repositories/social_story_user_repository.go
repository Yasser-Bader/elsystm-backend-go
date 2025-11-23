package repositories

import (
	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"gorm.io/gorm"
)

type SocialStoryUserRepository interface {
	Store(userId uint64, storyUserSeen dto.StoryUserSeen, storyId uint64) error
	Delete(userId uint64, storyId uint64) error
	UsersSeenStory(userId uint64, page int, storyId uint64, pagination *util.Pagination) error
}

type socialStoryUserRepository struct {
	db *gorm.DB
}

func NewSocialStoryUserRepository() SocialStoryUserRepository {
	return &socialStoryUserRepository{db: config.ConnectDB()}
}

func (r socialStoryUserRepository) Store(userId uint64, storyUserSeen dto.StoryUserSeen, storyId uint64) error {
	var storyUser []models.SocialStoryUser
	for _, id := range storyUserSeen.Ids {
		storyUser = append(storyUser, models.SocialStoryUser{
			UserID:  userId,
			StoryID: id,
		})
	}
	if err := r.db.Create(&storyUser).Error; err != nil {
		return err
	}
	return nil
}

func (r socialStoryUserRepository) Delete(userId uint64, storyId uint64) error {
	if err := r.db.Where("story_id=?", storyId).Delete(models.SocialStoryUser{}).Error; err != nil {
		return err
	}
	return nil
}

func (r socialStoryUserRepository) UsersSeenStory(userId uint64, page int, storyId uint64, pagination *util.Pagination) error {
	var users []models.SocialStoryUser
	var count int64
	if err := r.db.Model(&users).Where("user_id<>?", userId).Where("story_id=?", storyId).Count(&count).Error; err != nil {
		return err
	}
	err := r.db.Scopes(util.Paginate(page, pagination, count, 10)).Where("story_id=?", storyId).Preload("User").Preload("User.Avatar", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "App\\Domain\\Users\\Entity\\Models\\User").Where("collection_name=?", "avatar")
	}).Where("user_id<>?", userId).Order("id desc").Find(&users).Error
	if err != nil {
		return err
	}
	pagination.Data = users
	return nil
}

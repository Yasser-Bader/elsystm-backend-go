package repositories

import (
	"errors"
	"time"

	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SocialStoryRepository interface {
	Store(userId uint64, storyDTO dto.StoryCreateDTO, ctx *gin.Context) ([]models.SocialStoryCU, error)
	Share(userId uint64, id uint64) (models.SocialStoryCU, error)
	My(userId uint64) (models.UserWithStories, error)
	Archive(page int, userId uint64, pagination *util.Pagination) error
	Delete(userId uint64, id uint64) error
	User(id uint64, userId uint64) (models.UserWithStories, error)
	NewStories(userId uint64, followingId []uint64, page int, pagination *util.Pagination) error
	CountTodayStory(userId uint64) (int64, error)
}

type socialStoryRepository struct {
	db *gorm.DB
}

func NewStoryRepository() SocialStoryRepository {
	return &socialStoryRepository{db: config.ConnectDB()}
}

func (r socialStoryRepository) Store(userId uint64, storyDTO dto.StoryCreateDTO, ctx *gin.Context) ([]models.SocialStoryCU, error) {
	var stories []models.SocialStoryCU
	uploadedFiles := util.AWSUploadMultiableFiles(ctx, "files")
	if len(uploadedFiles) > 0 {
		for i, item := range uploadedFiles {
			var content string
			var duration string
			if len(storyDTO.Data) == len(uploadedFiles) {
				content = storyDTO.Data[i].Content
				duration = storyDTO.Data[i].Duration
			}
			stories = append(stories, models.SocialStoryCU{
				Content:  content,
				Duration: duration,
				UserId:   userId,
				FileName: item.URL,
				FileType: item.Filetype,
			})
		}
		if err := r.db.Create(&stories).Error; err != nil {
			return stories, err
		}
	}
	return stories, nil
}

func (r socialStoryRepository) Share(userId uint64, id uint64) (models.SocialStoryCU, error) {
	var story models.SocialStoryCU
	if err := r.db.Where("user_id=?", userId).Where("id=?", id).Find(&story).Error; err != nil {
		return story, err
	}
	if story.ID == 0 {
		return story, errors.New("not found")
	}
	story.CreatedAt = time.Now()
	story.UpdatedAt = time.Now()
	if err := r.db.Save(&story).Error; err != nil {
		return story, err
	}
	return story, nil
}

func (r socialStoryRepository) My(userId uint64) (models.UserWithStories, error) {
	var user models.UserWithStories
	var stories []models.SocialStoryWithoutSeenSerializer

	if err := r.db.Where("id=?", userId).Preload("Avatar", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "App\\Domain\\Users\\Entity\\Models\\User").Where("collection_name=?", "avatar")
	}).Find(&user).Error; err != nil {
		return user, err
	}

	err := r.db.Select("social_stories.*", "(select count(*) from `social_story_users` where `social_stories`.`id` = `social_story_users`.`story_id` and `social_stories`.`user_id` <> `social_story_users`.`user_id`) as `count_seen`").Preload("Seen", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id", userId)
	}).Where("social_stories.user_id=?", userId).Where("social_stories.created_at >= ?", time.Now().AddDate(0, 0, -1)).Order("social_stories.created_at desc").Find(&stories).Error
	if err != nil {
		return user, err
	}

	user.Stories = stories

	var isSeen bool = true
	for _, story := range user.Stories {
		if story.Seen.ID == 0 {
			isSeen = false
		}
	}
	user.SeenAll = isSeen
	return user, nil
}

func (r socialStoryRepository) Archive(page int, userId uint64, pagination *util.Pagination) error {
	var stories []models.SocialStory
	var count int64
	if err := r.db.Model(&stories).Where("user_id=?", userId).Where("created_at < ?", time.Now().AddDate(0, 0, -1)).Count(&count).Error; err != nil {
		return err
	}
	err := r.db.Scopes(util.Paginate(page, pagination, count, 5)).Select("social_stories.*", "(select count(*) from `social_story_users` where `social_stories`.`id` = `social_story_users`.`story_id`) as `count_seen`").Where("user_id=?", userId).Where("created_at < ?", time.Now().AddDate(0, 0, -1)).Order("created_at desc").Find(&stories).Error
	if err != nil {
		return err
	}
	pagination.Data = stories
	return nil
}

func (r socialStoryRepository) Delete(userId uint64, id uint64) error {
	var story models.SocialStory
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Delete(&story).Error; err != nil {
		return err
	}
	return nil
}

func (r socialStoryRepository) User(id uint64, userId uint64) (models.UserWithStories, error) {
	var user models.UserWithStories
	var stories []models.SocialStoryWithoutSeenSerializer

	// Get user
	if err := r.db.Where("id=?", id).Preload("Avatar", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "App\\Domain\\Users\\Entity\\Models\\User").Where("collection_name=?", "avatar")
	}).Find(&user).Error; err != nil {
		return user, err
	}

	// Get stories of user
	if err := r.db.Where("user_id=?", id).Preload("Seen", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id", userId)
	}).Where("created_at >= ?", time.Now().AddDate(0, 0, -1)).Order("created_at desc").Find(&stories).Error; err != nil {
		return user, err
	}

	user.Stories = stories

	var isSeen bool = true
	for _, story := range user.Stories {
		if story.Seen.ID == 0 {
			isSeen = false
		}
	}
	user.SeenAll = isSeen
	return user, nil
}

func (r socialStoryRepository) NewStories(userId uint64, followingId []uint64, page int, pagination *util.Pagination) error {
	var stories []models.SocialStoryWithoutSeenSerializer
	var users []models.UserWithStories
	var usersModel []models.UserWithStories
	var usersGroup []int64

	// Get my storie
	my, err := r.My(userId)
	if err != nil {
		return err
	}
	// Added my stories in list
	if len(my.Stories) > 0 {
		users = append(users, my)
	}

	// Check if new stories or no
	if err := r.db.Model(&stories).Where("user_id in(?)", followingId).Where("created_at >= ?", time.Now().AddDate(0, 0, -1)).Group("user_id").Pluck("user_id", &usersGroup).Error; err != nil {
		return err
	}
	if len(usersGroup) == 0 {
		pagination.Data = users
		return nil
	}

	// Get stories of following
	if err := r.db.Scopes(util.Paginate(page, pagination, int64(len(usersGroup)), 5)).Preload("Avatar", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "App\\Domain\\Users\\Entity\\Models\\User").Where("collection_name=?", "avatar")
	}).Preload("Stories", func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at >= ?", time.Now().AddDate(0, 0, -1)).Order("created_at desc")
	}).Preload("Stories.Seen", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userId)
	}).Where("id in(?)", usersGroup).Find(&usersModel).Error; err != nil {
		return err
	}

	// Added users in nodel
	users = append(users, usersModel...)

	// Check for seen
	for i, user := range users {
		var isSeen bool = true
		for _, story := range user.Stories {
			if story.Seen.ID == 0 {
				isSeen = false
			}
		}
		users[i].SeenAll = isSeen
	}

	// Update users
	pagination.Data = users
	return nil
}

func (r socialStoryRepository) CountTodayStory(userId uint64) (int64, error) {
	var count int64
	var story models.SocialStoryCU

	if err := r.db.Model(&story).Where("user_id=?", userId).Where("created_at >= ?", time.Now().AddDate(0, 0, -1)).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

package repositories

import (
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"gorm.io/gorm"
)

type HashtagRepository interface {
	Store(userId uint64, hashtagDTO dto.HashtagCreateDTO) (models.SocialHashtag, error)
	StoreByModel(userId uint64, postId uint64, hashtags []string, model string) error
	Delete(userId uint64, id uint64) error
	HashtagIsDuplicateName(name string) (ctx *gorm.DB)
}

type hashtagRepository struct {
	db *gorm.DB
}

func NewHashtagRepository() HashtagRepository {
	return &hashtagRepository{db: config.ConnectDB()}
}

func (r hashtagRepository) StoreByModel(userId uint64, postId uint64, hashtags []string, model string) error {
	var hashtag []models.SocialHashtagPCR

	for _, value := range hashtags {
		hashtagID, _ := strconv.ParseUint(value, 10, 64)
		if hashtagID != 0 {
			hashtag = append(hashtag, models.SocialHashtagPCR{
				HashtagId: hashtagID,
				ModelId:   postId,
				ModelType: model,
			})
		}
	}
	if err := r.db.Create(&hashtag).Error; err != nil {
		return err
	}

	return nil
}

func (r hashtagRepository) Store(userId uint64, hashtagDTO dto.HashtagCreateDTO) (models.SocialHashtag, error) {
	var hashtag models.SocialHashtag
	hashtag.UserId = userId
	hashtag.Name = hashtagDTO.Name
	if err := r.db.Create(&hashtag).Error; err != nil {
		return hashtag, err
	}
	return hashtag, nil
}

func (r hashtagRepository) Delete(userId uint64, id uint64) error {
	var hashtag models.SocialHashtag
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Delete(&hashtag).Error; err != nil {
		return err
	}
	return nil
}

func (repository *hashtagRepository) HashtagIsDuplicateName(name string) (ctx *gorm.DB) {
	var hashtag models.SocialHashtag
	return repository.db.Where("name = ?", name).Take(&hashtag)
}

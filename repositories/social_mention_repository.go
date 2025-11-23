package repositories

import (
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"gorm.io/gorm"
)

type MentionRepository interface {
	Store(userId uint64, postId uint64, mentions []string, model string) error
	Delete(userId uint64, id uint64) error
}

type mentionRepository struct {
	db *gorm.DB
}

func NewMentionRepository() MentionRepository {
	return &mentionRepository{db: config.ConnectDB()}
}

func (r mentionRepository) Store(userId uint64, postId uint64, mentions []string, model string) error {
	var postMentions []models.SocialMention
	if len(mentions) > 0 {
		for _, value := range mentions {
			mentionID, _ := strconv.ParseUint(value, 10, 64)
			if mentionID != 0 {
				postMentions = append(postMentions, models.SocialMention{
					UserId:        userId,
					UserMentionId: mentionID,
					ModelId:       postId,
					ModelType:     model,
				})
			}
		}
		if len(postMentions) > 0 {
			if err := r.db.Create(&postMentions).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (r mentionRepository) Delete(userId uint64, id uint64) error {
	var mention models.SocialMention
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Delete(&mention).Error; err != nil {
		return err
	}
	return nil
}

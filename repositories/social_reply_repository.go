package repositories

import (
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"gorm.io/gorm"
)

type ReplyRepository interface {
	Store(userId uint64, replyDTO dto.ReplyCreateDTO) (models.SocialReply, error)
	Delete(userId uint64, id uint64) error
	Update(userId uint64, id uint64, replyDTO dto.ReplyUpdateDTO) (models.SocialReply, error)
}

type replyRepository struct {
	db *gorm.DB
}

func NewReplyRepository() ReplyRepository {
	return &replyRepository{db: config.ConnectDB()}
}

func (r replyRepository) Store(userId uint64, replyDTO dto.ReplyCreateDTO) (models.SocialReply, error) {

	var reply models.SocialReply
	var replyHashtags []models.SocialHashtagPCR
	var replyMentions []models.SocialMention

	postId, _ := strconv.ParseUint(replyDTO.PostId, 10, 64)
	commentId, _ := strconv.ParseUint(replyDTO.CommentId, 10, 64)

	reply.PostId = postId
	reply.CommentId = commentId
	reply.Content = replyDTO.Content
	reply.UserId = userId
	if err := r.db.Create(&reply).Error; err != nil {
		return reply, err
	}

	// Set hashtags
	if len(replyDTO.Hashtags) > 0 {
		for _, value := range replyDTO.Hashtags {
			hashtagID, _ := strconv.ParseUint(value, 10, 64)
			if hashtagID != 0 {
				replyHashtags = append(replyHashtags, models.SocialHashtagPCR{
					HashtagId: hashtagID,
					ModelId:   reply.ID,
					ModelType: "reply",
				})
			}
		}
		if len(replyHashtags) > 0 {
			if err := r.db.Create(&replyHashtags).Error; err != nil {
				return reply, err
			}
		}
	}

	// Set mentions
	if len(replyDTO.Mentions) > 0 {
		for _, value := range replyDTO.Mentions {
			mentionID, _ := strconv.ParseUint(value, 10, 64)
			if mentionID != 0 {
				replyMentions = append(replyMentions, models.SocialMention{
					ModelId:   reply.ID,
					ModelType: "reply",
				})
			}
		}
		if len(replyMentions) > 0 {
			if err := r.db.Create(&replyMentions).Error; err != nil {
				return reply, err
			}
		}
	}

	return reply, nil
}

func (r replyRepository) Delete(userId uint64, id uint64) error {
	var reply models.SocialReply
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Delete(&reply).Error; err != nil {
		return err
	}
	return nil
}

func (r replyRepository) Update(userId uint64, id uint64, replyDTO dto.ReplyUpdateDTO) (models.SocialReply, error) {
	var reply models.SocialReply
	var history models.SocialHistory

	// Get reply
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Find(&reply).Error; err != nil {
		return reply, err
	}

	// Set reply
	history.UserId = userId
	history.Content = replyDTO.Content
	history.ModelId = id
	history.ModelType = "reply"

	if err := r.db.Create(&history).Error; err != nil {
		return reply, err
	}

	// Update Reply
	reply.Content = replyDTO.Content
	if err := r.db.Save(&reply).Error; err != nil {
		return reply, err
	}

	return reply, nil
}

package repositories

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Store(userId uint64, commentDTO dto.CommentCreateDTO) (models.SocialComment, error)
	Delete(userId uint64, id uint64) error
	Update(userId uint64, id uint64, commentDTO dto.CommentUpdateDTO) (models.SocialComment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository() CommentRepository {
	return &commentRepository{db: config.ConnectDB()}
}

func (r commentRepository) Store(userId uint64, commentDTO dto.CommentCreateDTO) (models.SocialComment, error) {
	var comment models.SocialComment
	var commentHashtags []models.SocialHashtagPCR
	var commentMentions []models.SocialMention

	postId, _ := strconv.ParseUint(commentDTO.PostId, 10, 64)
	comment.Content = commentDTO.Content
	comment.PostId = postId
	comment.UserId = userId
	if err := r.db.Create(&comment).Error; err != nil {
		return comment, err
	}

	// Set hashtags
	if len(commentDTO.Hashtags) > 0 {
		for _, value := range commentDTO.Hashtags {
			hashtagID, _ := strconv.ParseUint(value, 10, 64)
			if hashtagID != 0 {
				commentHashtags = append(commentHashtags, models.SocialHashtagPCR{
					HashtagId: hashtagID,
					ModelId:   comment.ID,
					ModelType: "comment",
				})
			}
		}
		if len(commentHashtags) > 0 {
			if err := r.db.Create(&commentHashtags).Error; err != nil {
				return comment, err
			}
		}
	}

	// Set mentions
	if len(commentDTO.Mentions) > 0 {
		for _, value := range commentDTO.Mentions {
			mentionID, _ := strconv.ParseUint(value, 10, 64)
			if mentionID != 0 {
				commentMentions = append(commentMentions, models.SocialMention{
					ModelId:   comment.ID,
					ModelType: "comment",
				})
			}
		}
		if len(commentMentions) > 0 {
			if err := r.db.Create(&commentMentions).Error; err != nil {
				return comment, err
			}
		}
	}

	return comment, nil
}

func (r commentRepository) Delete(userId uint64, id uint64) error {
	var comment models.SocialComment
	fmt.Println(userId)
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}

func (r commentRepository) Update(userId uint64, id uint64, commentDTO dto.CommentUpdateDTO) (models.SocialComment, error) {
	var comment models.SocialComment
	var history models.SocialHistory

	// Get comment
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Find(&comment).Error; err != nil {
		return comment, err
	}

	if comment.ID == 0 {
		return comment, errors.New("not found")
	}

	// Set History
	history.UserId = userId
	history.Content = commentDTO.Content
	history.ModelId = id
	history.ModelType = "comment"

	if err := r.db.Create(&history).Error; err != nil {
		return comment, err
	}

	// Update Comment
	comment.Content = commentDTO.Content
	if err := r.db.Save(&comment).Error; err != nil {
		return comment, err
	}

	return comment, nil
}

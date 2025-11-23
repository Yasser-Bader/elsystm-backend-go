package repositories

import (
	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"gorm.io/gorm"
)

type LikeRepository interface {
	Store(userId uint64, likeDTO dto.LikeCreateDTO) (models.SocialLike, error)
	Delete(userId uint64, id uint64) error
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository() LikeRepository {
	return &likeRepository{db: config.ConnectDB()}
}

func (r likeRepository) Store(userId uint64, likeDTO dto.LikeCreateDTO) (models.SocialLike, error) {
	var like models.SocialLike

	like.UserId = userId
	like.ModelId = likeDTO.ModelId
	like.ModelType = likeDTO.ModelType
	if err := r.db.Where("user_id=?", userId).Where("model_type=?", likeDTO.ModelType).Where("model_id=?", likeDTO.ModelId).FirstOrCreate(&like).Error; err != nil {
		return like, err
	}
	return like, nil
}

func (r likeRepository) Delete(userId uint64, id uint64) error {
	var like models.SocialLike
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Delete(&like).Error; err != nil {
		return err
	}
	return nil
}

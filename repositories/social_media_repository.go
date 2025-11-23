package repositories

import (
	"net/http"

	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SocialMediaRepository interface {
	Store(userId uint64, modelId uint64, modelType string, ctx *gin.Context) ([]models.SocialMedia, error)
	StoreSingleFile(userId uint64, modelId uint64, modelType string, ctx *gin.Context) (models.SocialMedia, error)
	StoreUpdate(userId uint64, modelId uint64, modelType string, ctx *gin.Context) ([]models.SocialMedia, error)
}

type socialMediaRepository struct {
	db *gorm.DB
}

func NewSocialMediaRepository() SocialMediaRepository {
	return &socialMediaRepository{db: config.ConnectDB()}
}

func (r socialMediaRepository) Store(userId uint64, modelId uint64, modelType string, ctx *gin.Context) ([]models.SocialMedia, error) {
	var media []models.SocialMedia
	uploadedFiles := util.AWSUploadMultiableFiles(ctx, "media")
	if len(uploadedFiles) > 0 {
		for _, item := range uploadedFiles {
			media = append(media, models.SocialMedia{
				UserId:    userId,
				ModelType: modelType,
				ModelId:   modelId,
				FileName:  item.URL,
				FileType:  item.Filetype,
			})
		}
		if err := r.db.Create(&media).Error; err != nil {
			return media, err
		}
	}

	return media, nil
}

func (r socialMediaRepository) StoreSingleFile(userId uint64, modelId uint64, modelType string, ctx *gin.Context) (models.SocialMedia, error) {
	var socialMedia models.SocialMedia
	file, _, _ := ctx.Request.FormFile("file")
	if file != nil {
		filename, fileType, err := util.AWSUploadSingleFile(ctx, "file")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": "Error"})
			return socialMedia, err
		}
		socialMedia.ModelId = modelId
		socialMedia.ModelType = modelType
		socialMedia.FileName = filename
		socialMedia.FileType = fileType
		socialMedia.UserId = userId
		r.db.Create(&socialMedia)
	}
	return socialMedia, nil
}

func (r socialMediaRepository) StoreUpdate(userId uint64, modelId uint64, modelType string, ctx *gin.Context) ([]models.SocialMedia, error) {
	var media []models.SocialMedia
	// Get post
	if err := r.db.Where("model_id=?", modelId).Where("user_id=?", userId).Where("model_type=post").Find(&media).Error; err != nil {
		return media, err
	}
	uploadedFiles := util.AWSUploadMultiableFiles(ctx, "media")
	if len(uploadedFiles) > 0 {
		for _, item := range uploadedFiles {
			media = append(media, models.SocialMedia{
				UserId:    userId,
				ModelType: modelType,
				ModelId:   modelId,
				FileName:  item.URL,
				FileType:  item.Filetype,
			})
		}
		if err := r.db.Save(&media).Error; err != nil {
			return media, err
		}
	}

	return media, nil
}

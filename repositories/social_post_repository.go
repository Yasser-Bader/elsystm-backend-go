package repositories

import (
	"github.com/Elsystm-Inc/systm-go-social/config"
	"github.com/Elsystm-Inc/systm-go-social/dto"
	"github.com/Elsystm-Inc/systm-go-social/models"
	"github.com/Elsystm-Inc/systm-go-social/util"
	"gorm.io/gorm"
)

type PostRepository interface {
	Index(userId uint64, page int, pagination *util.Pagination) error
	Store(userId uint64, postDTO dto.PostCreateDTO) (models.SocialPost, error)
	Delete(userId uint64, id uint64) error
	Update(userId uint64, id uint64, postDTO dto.PostUpdateDTO) (models.SocialPost, error)
	GetPost(id uint64) (models.SocialGetPosts, error)
}
type postRepository struct {
	db *gorm.DB
}

func NewPostRepository() PostRepository {
	return &postRepository{db: config.ConnectDB()}
}
func (r postRepository) Index(userId uint64, page int, pagination *util.Pagination) error {
	var posts []models.SocialPostHome
	if err := r.db.Select("social_posts.*",
		"(select count(*) from `social_comments` where `social_posts`.`id` = `social_comments`.`post_id`) as `count_comment`",
		"(select count(*) from `social_likes` where `social_posts`.`id` = `social_likes`.`model_id` and `model_type`='post') as `count_like`",
		"(select count(*) from `social_post_shares` where `social_posts`.`id` = `social_post_shares`.`post_share_id`) as `count_share`").Scopes(util.Paginate(page, pagination, 0, 5)).Preload("Media", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "post")
	}).Preload("User").Preload("User.Avatar", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "App\\Domain\\Users\\Entity\\Models\\User").Where("collection_name=?", "avatar")
	}).Preload("SharedPost").Preload("SharedPost.PostShared").Preload("SharedPost.PostShared.Media", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "post")
	}).Preload("SharedPost.PostShared.User").Preload("SharedPost.PostShared.User.Avatar", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "App\\Domain\\Users\\Entity\\Models\\User").Where("collection_name=?", "avatar")
	}).Preload("Mentions", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type='post'")
	}).Preload("Hashtags", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type='post'")
	}).Preload("Hashtags.Hashtag").Preload("Mentions.User").Order("updated_at desc").Find(&posts).Error; err != nil {
		return err
	}
	pagination.Data = posts
	return nil
}

func (r postRepository) Store(userId uint64, postDTO dto.PostCreateDTO) (models.SocialPost, error) {
	var post models.SocialPost
	var postShare models.SocialPostShare

	post.Content = postDTO.Content
	post.UserId = userId
	if err := r.db.Create(&post).Error; err != nil {
		return post, err
	}

	// Set Share
	if postDTO.SharedPostId != 0 {
		postShare.PostId = post.ID
		postShare.PostShareId = postDTO.SharedPostId
		if err := r.db.Create(&postShare).Error; err != nil {
			return post, err
		}
	}

	return post, nil
}

func (r postRepository) Delete(userId uint64, id uint64) error {
	var post models.SocialPost
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Delete(&post).Error; err != nil {
		return err
	}
	return nil
}

func (r postRepository) Update(userId uint64, id uint64, postDTO dto.PostUpdateDTO) (models.SocialPost, error) {
	var post models.SocialPost
	var history models.SocialHistory

	// Get post
	if err := r.db.Where("id=?", id).Where("user_id=?", userId).Find(&post).Error; err != nil {
		return post, err
	}

	// Set post
	history.UserId = userId
	history.Content = postDTO.Content
	history.ModelId = id
	history.ModelType = "post"

	if err := r.db.Create(&history).Error; err != nil {
		return post, err
	}
	// Update Post
	post.Content = postDTO.Content
	if err := r.db.Save(&post).Error; err != nil {
		return post, err
	}

	return post, nil
}

func (r postRepository) GetPost(id uint64) (models.SocialGetPosts, error) {
	var post []models.SocialPost
	var getPost models.SocialGetPosts
	var comments []models.SocialComment
	var replies []models.SocialReply
	var likes []models.SocialLike
	var countComments int64
	var countReply int64
	var countPostLikes int64
	var commentIds []uint64
	if err := r.db.Model(&comments).Where("post_id=?", id).Count(&countComments).Pluck("id", &commentIds).Error; err != nil {
		return getPost, err
	}
	if err := r.db.Model(&replies).Where("post_id=?", id).Count(&countReply).Error; err != nil {
		return getPost, err
	}
	if err := r.db.Model(&likes).Where("model_id=?", id).Where("model_type=?", "post").Count(&countPostLikes).Error; err != nil {
		return getPost, err
	}

	// posts
	err := r.db.Debug().Preload("Media", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id=?", id).Where("model_type=?", "post")
	}).Preload("Hashtags", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id=?", id).Where("model_type=?", "post")
	}).Preload("Mentions", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id=?", id).Where("model_type=?", "post")

		// comments
	}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Where("post_id=?", id)
	}).Preload("Comments.Hashtags", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id IN ?", commentIds).Where("model_type=?", "comment")
	}).Preload("Comments.Mentions", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id IN ?", commentIds).Where("model_type=?", "comment")
	}).Preload("Comments.Likes", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id IN ?", commentIds).Where("model_type=?", "comment")
	}).Preload("Comments").Preload("Comments.Media", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id IN ?", commentIds).Where("model_type=?", "comment")

		// replies
	}).Preload("Comments").Preload("Comments.Replies", func(db *gorm.DB) *gorm.DB {
		return db.Where("post_id=?", id).Where("comment_id IN ?", commentIds)
	}).Preload("Comments").Preload("Comments.Replies.Media", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_type=?", "reply")
	}).Preload("Comments").Preload("Comments.Replies.Hashtags", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id IN ?", commentIds).Where("model_type=?", "reply")
	}).Preload("Comments").Preload("Comments.Replies.Mentions", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id IN ?", commentIds).Where("model_type=?", "reply")
	}).Preload("Comments").Preload("Comments.Replies.Likes", func(db *gorm.DB) *gorm.DB {
		return db.Where("model_id IN ?", commentIds).Where("model_type=?", "reply")

	}).Where("id=?", id).Find(&post).Error
	if err != nil {
		return getPost, err
	}

	getPost.Data = post
	getPost.Comments = countComments + countReply
	getPost.Likes = countPostLikes
	return getPost, nil
}

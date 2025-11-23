package dto

type HashtagCreateDTO struct {
	Name string `form:"name" binding:"required"`
}

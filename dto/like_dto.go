package dto

type LikeCreateDTO struct {
	ModelId   uint64 `json:"model_id" binding:"required"`
	ModelType string `json:"model_type" binding:"required"`
}

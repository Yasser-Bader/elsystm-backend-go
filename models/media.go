package models

type Media struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	ModelId  uint64 `json:"model_id"`
	Name     string `json:"name"`
	FileName string `json:"file_name"`
}

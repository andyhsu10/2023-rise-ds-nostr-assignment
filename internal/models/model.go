package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
}

type CoreEvent struct {
	BaseModel
	Name string `json:"name"`
	Data string `json:"data"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}

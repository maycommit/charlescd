package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
// 	b.ID = uuid.New()
// 	return nil
// }

func (u *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// func (BaseModel *BaseModel) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("ID", uuid.New())
// 	return nil
// }

package base

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel holds the fields shared by all business models.
type BaseModel struct {
	// ID is the primary key.
	ID int64 `gorm:"primarykey;type:bigint" json:"id"`
	// CreatedAt records the creation time.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt records the last update time.
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt records the soft delete time.
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

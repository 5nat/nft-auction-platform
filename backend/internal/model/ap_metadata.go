package model

import "time"

type AppMetadata struct {
	MetaKey   string    `gorm:"column:meta_key;primaryKey;type:varchar(128)" json:"key"`
	Value     string    `gorm:"type:varchar(255);not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AppMetadata) TableName() string {
	return "app_metadata"
}

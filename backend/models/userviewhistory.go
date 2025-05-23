package models

import (
	"time"

	"gorm.io/gorm"
)

// UserViewHistory tracks anime details pages viewed by a user.
type UserViewHistory struct {
	gorm.Model
	UserID          uint      `json:"user_id" gorm:"not null;index:idx_user_anime_view,unique"`
	AnimeExternalID int       `json:"anime_external_id" gorm:"not null;index:idx_user_anime_view,unique"` // AniList ID
	LastViewedAt    time.Time `json:"last_viewed_at" gorm:"index"`
	ViewCount       uint      `json:"view_count" gorm:"default:1"`

	// User User `gorm:"foreignKey:UserID"` // Optional: for GORM relations if needed directly
}

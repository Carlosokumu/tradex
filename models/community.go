package models

import (
	"time"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name        string      `gorm:"size:100;not null"`
	Description string      `gorm:"size:500"`
	AdminID     uint        `gorm:"not null"`
	Admin       UserModel   `gorm:"foreignKey:AdminID"`
	Members     []UserModel `gorm:"many2many:user_communities;"`
	Posts       []Post      `gorm:"foreignKey:CommunityID"`
}

type Post struct {
	gorm.Model
	Content     string    `gorm:"size:1000"`
	ImageURL    string    `gorm:"not null"`
	PosterID    uint      `gorm:"not null"`
	Poster      UserModel `gorm:"foreignKey:PosterID"`
	CommunityID uint      `gorm:"not null"`
	Community   Community `gorm:"foreignKey:CommunityID"`
	Comments    []Comment `gorm:"foreignKey:PostID"`
	Likes       []Like    `gorm:"foreignKey:PostID"`
}

type Comment struct {
	gorm.Model
	Content     string    `gorm:"size:500;not null"`
	CommenterID uint      `gorm:"not null"`
	Commenter   UserModel `gorm:"foreignKey:CommenterID"`
	PostID      uint      `gorm:"not null"`
	Post        Post      `gorm:"foreignKey:PostID"`
}

type Like struct {
	gorm.Model
	UserID uint      `gorm:"not null"`
	User   UserModel `gorm:"foreignKey:UserID"`
	PostID uint      `gorm:"not null"`
	Post   Post      `gorm:"foreignKey:PostID"`
}

// Explicit join table model to control the relationship
type CommunityMember struct {
	UserID      uint      `gorm:"primaryKey;column:user_id"`
	CommunityID uint      `gorm:"primaryKey;column:community_id"`
	JoinedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	User        UserModel `gorm:"foreignKey:UserID"`
	Community   Community `gorm:"foreignKey:CommunityID"`
}

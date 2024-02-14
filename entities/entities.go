package entities

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email         string `gorm:"unique;not null;check:LENGTH(email) >= 3 AND email LIKE '%@%'"`
	Username      string `gorm:"unique;not null;check:LENGTH(username) > 0"`
	PasswordHash  string
	Administrator bool `gorm:"not null;default false"`
	Banned        bool `gorm:"not null;default false"`
	Reviews       []PlaygroundReview
	Photos        []PlaygroundPhoto
	ReviewVotes   []ReviewVote
	PhotoVotes    []PhotoVote
}

type Playground struct {
	gorm.Model
	SiteNumber string `gorm:"unique;not null"`
	//site number can be used as identification of the site when reporting an issue to the authorities for example
	Latitude       float64
	Longitude      float64
	Area           int `gorm:"check:area >= 0"`
	Location       string
	Ownership      string
	Reviews        []PlaygroundReview
	Photos         []PlaygroundPhoto
	SelectedPhotos []PlaygroundPhoto
	AverageRating  float32 `gorm:"-"`
}

//TODO: add relation between playgrounds and their selected photos to be shown

type PlaygroundReview struct {
	gorm.Model
	Playground   Playground
	User         User
	PlaygroundID uint `gorm:"not null"`
	UserID       uint `gorm:"not null"`
	Stars        int  `gorm:"check:stars >= 0 AND stars <= 5"`
	Content      string
	Votes        []ReviewVote
}

type PlaygroundPhoto struct {
	gorm.Model
	Playground   Playground
	User         User
	PlaygroundID uint   `gorm:"not null"`
	UserId       uint   `gorm:"not null"`
	Extension    string `gorm:"not null"`
	Approved     *bool
	//NULL is not reviewed yet, TRUE is approved and FALSE is rejected
	Selected bool `gorm:"not null;default false"`
	Votes    []PhotoVote
}

type ReviewVote struct {
	Up                 bool
	Review             PlaygroundReview `gorm:"foreignkey:PlaygroundReviewID"`
	User               User             `gorm:"foreignkey:UserID"`
	PlaygroundReviewID uint             `gorm:"primaryKey;uniqueIndex:idx_review_usr_vote"`
	UserID             uint             `gorm:"primaryKey;uniqueIndex:idx_review_usr_vote"`
	CreatedAt          time.Time        `gorm:"<-:create"`
	UpdatedAt          time.Time
	gorm.DeletedAt
}

type PhotoVote struct {
	Up                bool
	Photo             PlaygroundPhoto `gorm:"foreignkey:PlaygroundPhotoID"`
	User              User            `gorm:"foreignkey:UserID"`
	PlaygroundPhotoID uint            `gorm:"primaryKey;uniqueIndex:idx_photo_usr_vote"`
	UserID            uint            `gorm:"primaryKey;uniqueIndex:idx_photo_usr_vote"`
	CreatedAt         time.Time       `gorm:"<-:create"`
	UpdatedAt         time.Time
	gorm.DeletedAt
}

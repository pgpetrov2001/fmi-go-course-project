package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string `gorm:"unique;not null"`
	Username      string `gorm:"unique;not null"`
	PasswordHash  string
	Administrator bool `gorm:"not null;default false"`
	Banned        bool `gorm:"not null;default false"`
}

type Playground struct {
	gorm.Model
	SiteNumber string `gorm:"unique;not null"`
	//site number can be used as identification of the site when reporting an issue to the authorities for example
	Latitude  float64
	Longitude float64
	Area      int `gorm:"check:area >= 0"`
	Location  string
	Ownership string
}

//TODO: add relation between playgrounds and their selected photos to be shown

type PlaygroundReview struct {
	gorm.Model
	Playground Playground
	User       User
	Stars      int `gorm:"check:stars >= 0 AND stars <= 5"`
	Content    string
	Upvotes    int `gorm:"check:upvotes >= 0"`
	Downvotes  int `gorm:"check:downvotes <= 0"`
}

type PlaygroundPhoto struct {
	gorm.Model
	Playground Playground
	User       User
	Approved   bool
	Selected   bool
	//NULL is not reviewed yet, TRUE is approved and FALSE is rejected
	Upvotes   int `gorm:"check:upvotes >= 0"`
	Downvotes int `gorm:"check:downvotes <= 0"`
}

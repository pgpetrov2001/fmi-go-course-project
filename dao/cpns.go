package dao

import (
	"course-project/entities"
	"course-project/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path"
	"path/filepath"
)

type CPNS struct {
	Db            *gorm.DB
	FSStoragePath string
}

func (cpns *CPNS) Init() error {
	err := os.MkdirAll(cpns.FSStoragePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Could not intialize CPNS: %v", err)
	}
	return nil
}

func (cpns *CPNS) CreateUser(email string, username string, pass string, administrator bool, banned bool) (*entities.User, error) {
	hash, err := utils.HashPassword(pass)
	if err != nil {
		panic(err)
	}
	user := entities.User{
		Email:         email,
		Username:      username,
		PasswordHash:  hash,
		Administrator: administrator,
		Banned:        banned,
	}
	err = cpns.Db.Create(&user).Error
	return &user, err
}

func (cpns *CPNS) UpdateUser(user *entities.User) error {
	result := cpns.Db.Save(user)
	if result.RowsAffected == 0 {
		return fmt.Errorf("Could not find user with id %d, did not update anything", user.ID)
	}
	return result.Error
}

func (cpns *CPNS) DeleteUser(user *entities.User) error {
	result := cpns.Db.Delete(user)
	if result.RowsAffected == 0 {
		return fmt.Errorf("No user with id %d, did not delete anything", user.ID)
	}
	return result.Error
}

func (cpns *CPNS) Authenticate(email string, pass string) (*entities.User, error) {
	var user entities.User
	err := cpns.Db.First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if !utils.CheckPasswordHash(pass, user.PasswordHash) {
		return nil, nil
	}
	return &user, nil
}

func (cpns *CPNS) GetUsers() ([]entities.User, error) {
	var users []entities.User
	err := cpns.Db.Find(&users).Error
	return users, err
}

func (cpns *CPNS) GetUser(userId uint) (entities.User, error) {
	var user entities.User
	err := cpns.Db.First(&user, userId).Error
	return user, err
}

func (cpns *CPNS) GetPlayground(playgroundId uint) (entities.Playground, error) {
	var playground entities.Playground
	err := cpns.Db.Preload("Photos").Preload("Reviews").First(&playground, playgroundId).Error
	return playground, err
}

func (cpns *CPNS) GetPlaygrounds() ([]entities.Playground, error) {
	var playgrounds []entities.Playground
	err := cpns.Db.Preload("Photos").Preload("Reviews").Find(&playgrounds).Error
	return playgrounds, err
}

func (cpns *CPNS) CreatePlayground(playground *entities.Playground) error {
	result := cpns.Db.Create(playground)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		return fmt.Errorf("Could not create playground entry in database for some reason")
	}
	return err
}

func (cpns *CPNS) UpdatePlayground(playground *entities.Playground) error {
	result := cpns.Db.Save(playground)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		return fmt.Errorf("Could not find playground with id %d, did not update anything", playground.ID)
	}
	return result.Error
}

func (cpns *CPNS) DeletePlayground(playground *entities.Playground) error {
	result := cpns.Db.Delete(playground)
	if result.RowsAffected == 0 {
		return fmt.Errorf("No playground with id %d, did not delete anything", playground.ID)
	}
	return result.Error
}

func (cpns *CPNS) ReviewPlayground(review *entities.PlaygroundReview) error {
	alreadyRated := false
	for _, currentReview := range review.Playground.Reviews {
		if currentReview.UserID == review.UserID {
			alreadyRated = true
			break
		}
	}
	if alreadyRated {
		return fmt.Errorf("User %s has already submitted a review for playground with id %u", review.User.Username, review.PlaygroundID)
	}
	result := cpns.Db.Create(review)
	err := result.Error
	if err != nil && result.RowsAffected == 0 {
		return fmt.Errorf("Somehow the review couldn't be created")
	}
	return err
}

func (cpns *CPNS) PlaygroundGallery(playgroundId uint) ([]entities.PlaygroundPhoto, error) {
	var photos []entities.PlaygroundPhoto
	err := cpns.Db.Where("playground_id = ?", playgroundId).Find(&photos).Error
	return photos, err
}

func (cpns *CPNS) PendingPhotos() ([]entities.PlaygroundPhoto, error) {
	var photos []entities.PlaygroundPhoto
	err := cpns.Db.Where("approved IS NULL").Find(&photos).Error
	return photos, err
}

func (cpns *CPNS) UploadPhoto(photo *entities.PlaygroundPhoto, filename string, data []byte) error {
	photo.Extension = filepath.Ext(filename)
	result := cpns.Db.Create(photo)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		return fmt.Errorf("For some reason the database could not create the specified photo")
	}
	if err != nil {
		return err
	}
	path := filepath.Join(cpns.FSStoragePath, "photos", fmt.Sprintf("%u.%s", photo.ID, photo.Extension))
	err = os.WriteFile(path, data, 0644)
	return err
}

func (cpns *CPNS) GetReview(reviewId uint) (entities.PlaygroundReview, error) {
	var review entities.PlaygroundReview
	err := cpns.Db.First(&review, reviewId).Error
	return review, err
}

func (cpns *CPNS) UpdateReview(review *entities.PlaygroundReview) error {
	result := cpns.Db.Save(review)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		return fmt.Errorf("Could not find review with id %d, did not update anything", review.ID)
	}
	return result.Error
}

func (cpns *CPNS) GetPhoto(photoId uint) (entities.PlaygroundPhoto, error) {
	var photo entities.PlaygroundPhoto
	err := cpns.Db.First(&photo, photoId).Error
	return photo, err
}

func (cpns *CPNS) GetPhotoContents(photo *entities.PlaygroundPhoto) ([]byte, error) {
	fillePath := path.Join(cpns.FSStoragePath, "photos", fmt.Sprintf("%u.%s", photo.ID, photo.Extension))
	return os.ReadFile(fillePath)
}

func (cpns *CPNS) UpdatePhoto(photo *entities.PlaygroundPhoto) error {
	result := cpns.Db.Save(photo)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		return fmt.Errorf("Could not find photo with id %d, did not update anything", photo.ID)
	}
	return result.Error
}

func (cpns *CPNS) VoteReview(review *entities.ReviewVote) error {
	result := cpns.Db.Create(review)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		return fmt.Errorf("For some reason the database could not create a vote")
	}
	return err
}

func (cpns *CPNS) VotePhoto(photo *entities.PhotoVote) error {
	result := cpns.Db.Create(photo)
	err := result.Error
	if err == nil && result.RowsAffected == 0 {
		return fmt.Errorf("For some reason the database could not create a vote")
	}
	return err
}

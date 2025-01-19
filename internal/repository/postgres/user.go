package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time
	Avatar    string
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"uniqueIndex;not null"`
	Secret    string         `gorm:"uniqueIndex;not null"`
}

func (r *Repository) DeleteUser(ctx context.Context, secret string) (bool, error) {
	if err := r.db.Where("secret = ?", secret).Delete(&User{}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) GetUser(ctx context.Context, secret string) (*User, error) {
	var user User
	if err := r.db.Where("secret = ?", secret).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserBySecret(secret string) (*User, error) {
	var user User
	if err := r.db.Where("secret = ?", secret).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

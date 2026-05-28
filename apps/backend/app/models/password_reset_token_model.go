package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PasswordResetToken struct {
	BaseModelWithoutUpdatedAt
	Email     string    `gorm:"not null;size:100" faker:"email"`
	Token     string    `gorm:"not null;size:255" faker:"password"`
	ExpiredAt time.Time `gorm:"not null" faker:"date"`
	UsedAt    *time.Time
}

// HashToken mengenkripsi token sebelum disimpan
func (p *PasswordResetToken) HashToken(token *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*token), 14)
	if err != nil {
		return err
	}
	p.Token = string(bytes)
	return nil
}

// CheckToken memverifikasi token
func (p *PasswordResetToken) CheckToken(providedToken string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.Token), []byte(providedToken))
}

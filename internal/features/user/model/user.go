package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Username     string    `gorm:"type:varchar(100);not null" json:"username"`
	Email        string    `gorm:"type:varchar(255);unique;index" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255)" json:"-"`
	AuthProvider string    `gorm:"type:varchar(50);default:'EMAIL'" json:"auth_provider"` // EMAIL, GOOGLE, GUEST
	FCMToken     string    `gorm:"type:text" json:"fcm_token"`
	Timezone     string    `gorm:"type:varchar(50);default:'Asia/Jakarta'" json:"timezone"`
	LastActiveAt time.Time `json:"last_active_at"`

	// Audit Trail
	CreatedBy string `gorm:"type:varchar(36)" json:"created_by"`
	UpdatedBy string `gorm:"type:varchar(36)" json:"updated_by"`
	DeletedBy string `gorm:"type:varchar(36)" json:"deleted_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

// SetPassword encrypts the password using bcrypt
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword compares the hashed password with plain text
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

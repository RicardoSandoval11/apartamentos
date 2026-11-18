package entities

import (
	"time"

	"gorm.io/gorm"
)

const (
	IdentityCreated = iota
	IdentityVerified
	IdentityRegistered
	IdentityInactive
	IdentityBlocked
)

type Identity struct {
	Id           uint64     `gorm:"column:id;primaryKey;autoIncremental"`
	Email        string     `gorm:"column:email;size:75;not null;unique;index:idx_email"`
	PasswordHash string     `gorm:"column:password_hash;size:250;not null"`
	Status       int        `gorm:"column:status;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	RoleId       uint64     `gorm:"column:role_id;not null"`
	Role         Role       `gorm:"foreignKey:RoleId;references:Id;constraint:OnDelete:Cascade"`
}

func (identity *Identity) BeforeCreate(tx *gorm.DB) error {
	identity.CreatedAt = time.Now()
	return nil
}

func (identity *Identity) BeforeUpdate(tx *gorm.DB) error {
	currentTime := time.Now()
	identity.UpdatedAt = &currentTime
	return nil
}

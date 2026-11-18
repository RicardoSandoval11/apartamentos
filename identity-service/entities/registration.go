package entities

import (
	"time"

	"gorm.io/gorm"
)

type Registration struct {
	FirstName   string     `gorm:"column:first_name;size:50;not null"`
	LastName    string     `gorm:"column:last_name;size:70;not null"`
	PhoneNumber int        `gorm:"column:hone_number;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (registration *Registration) BeforeCreate(tx *gorm.DB) error {
	registration.CreatedAt = time.Now()
	return nil
}

func (registration *Registration) BeforeUpdate(tx *gorm.DB) error {
	time := time.Now()
	registration.UpdatedAt = &time
	return nil
}

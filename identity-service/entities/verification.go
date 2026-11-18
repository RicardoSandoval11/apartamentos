package entities

import "time"

type Verification struct {
	Id         uint64    `gorm:"column:id;primaryKey;autoIncremental"`
	Code       string    `gorm:"column:code;size:8;not null"`
	Expiration time.Time `gorm:"column:expiration;not null"`
	IdentityId uint64    `gorm:"column:identity_id;not null"`
	Identity   Identity  `gorm:"foreignKey:IdentityId;references:Id;constraint:OnDelete:Cascade"`
	Verified   bool      `gorm:"column:verified;not null"`
}

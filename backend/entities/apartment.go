package entities

import "github.com/google/uuid"

type Apartment struct {
	Id       uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	PublicId uuid.UUID `gorm:"column:public_id;not null;unique"`
	Title    string    `gorm:"column:title;not null;size:150"`
}

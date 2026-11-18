package entities

type Role struct {
	Id   uint64 `gorm:"column:id;primaryKey;autoIncremental"`
	Name string `gorm:"column:name;size:50;not null;unique"`
}

package entities

type Permission struct {
	Id     uint64 `gorm:"column:id;primaryKey;autoIncremental"`
	Name   string `gorm:"column:name;size:150;not null;uniqueIndex:uq_name_role"`
	RoleId uint64 `gorm:"column:role_id;not null;index:role_id;uniqueIndex:uq_name_role"`
	Role   Role   `gorm:"foreignKey:RoleId;references:Id;constraint:OnDelete:Cascade"`
}

package entities

type Attempt struct {
	Id             uint64       `gorm:"column:id;primaryKey;autoIncremental"`
	Times          int          `gorm:"column:times;not null"`
	VerificationId uint64       `gorm:"column:verification_id;not null"`
	Verification   Verification `gorm:"foreignKey:VerificationID;references:Id;constraint:OnDelete:Cascade"`
}

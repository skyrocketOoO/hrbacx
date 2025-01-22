package model

type Edge struct {
	ID uint `gorm:"primarykey"`

	From     string `gorm:"type:varchar(36)"`
	Relation string `gorm:"type:varchar(255)"`
	To       string `gorm:"type:varchar(36)"`
}

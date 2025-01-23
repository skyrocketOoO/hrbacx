package model

type Edge struct {
	ID uint `gorm:"primarykey"`

	From     string `gorm:"type:varchar(36);uniqueIndex:idx_from_relation_to"`
	Relation string `gorm:"type:varchar(255);uniqueIndex:idx_from_relation_to"`
	To       string `gorm:"type:varchar(36);uniqueIndex:idx_from_relation_to"`
}

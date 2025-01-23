package model

type Edge struct {
	ID uint `gorm:"primarykey"`

	FromV    string `gorm:"type:varchar(36);uniqueIndex:idx_from_relation_to"`
	Relation string `gorm:"type:varchar(255);uniqueIndex:idx_from_relation_to"`
	ToV      string `gorm:"type:varchar(36);uniqueIndex:idx_from_relation_to"`
}

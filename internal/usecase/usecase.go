package usecase

import (
	"gorm.io/gorm"
)

type Usecase struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Usecase {
	return &Usecase{db}
}

package global

import (
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

const (
	ApiVersion = "v1"

	// relation
	LearderOf = "leader_of"
	BelongsTo = "belongs_to"
)

var (
	// env | flag
	Database    string
	AutoMigrate bool = false
	Env         string

	// instance
	DB        *gorm.DB
	Validator *validator.Validate // use a single instance of Validate, it caches struct info
)

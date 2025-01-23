package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/hrbacx/internal/global"
)

func New() {
	log.Info().Msg("InitValidator")
	global.Validator = validator.New(validator.WithRequiredStructEnabled())
}

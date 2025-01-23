package boot

import (
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/hrbacx/internal/service/database"
	"github.com/skyrocketOoO/hrbacx/internal/service/inter/validator"
)

func NewService() error {
	log.Info().Msg("InitService")
	if err := database.New(); err != nil {
		return err
	}

	validator.New()

	return nil
}

package boot

import (
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/hrbacx/internal/global"
	"github.com/skyrocketOoO/hrbacx/internal/service/exter/database"
	nebulaservice "github.com/skyrocketOoO/hrbacx/internal/service/exter/nebula"
	"github.com/skyrocketOoO/hrbacx/internal/service/inter/validator"
)

func NewService() error {
	log.Info().Msg("InitService")

	if global.Database == "nebula" {
		if err := nebulaservice.New(); err != nil {
			return err
		}
	} else if global.Database == "postgres" {
		if err := database.New(); err != nil {
			return err
		}
	}

	validator.New()

	return nil
}

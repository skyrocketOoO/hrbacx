package boot

import (
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/hrbacx/internal/global"
	"github.com/skyrocketOoO/hrbacx/internal/service/exter/database"
	nebulaservice "github.com/skyrocketOoO/hrbacx/internal/service/exter/nebula"
	redisservice "github.com/skyrocketOoO/hrbacx/internal/service/exter/redis"
	"github.com/skyrocketOoO/hrbacx/internal/service/inter/validator"
)

func NewService() error {
	log.Info().Msg("InitService")

	if global.Database == "nebula" {
		if err := nebulaservice.New(); err != nil {
			return err
		}
	} else if global.Database == "redis" {
		if err := redisservice.New(); err != nil {
			return err
		}
	} else {
		if err := database.New(); err != nil {
			return err
		}
	}

	validator.New()

	return nil
}

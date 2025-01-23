package boot

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitConfig() (err error) {
	log.Info().Msg("InitConfig")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".env")

	if err = viper.ReadInConfig(); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

package server

import (
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/hrbacx/api"
	"github.com/skyrocketOoO/hrbacx/internal/boot"
	"github.com/skyrocketOoO/hrbacx/internal/controller"
	"github.com/skyrocketOoO/hrbacx/internal/global"
	"github.com/skyrocketOoO/hrbacx/internal/usecase"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "The main service command",
	Long:  ``,
	// Args:  cobra.MinimumNArgs(1),
	Run: RunServer,
}

func RunServer(cmd *cobra.Command, args []string) {
	if err := boot.InitAll(); err != nil {
		log.Fatal().Msgf("Initialization failed: %v", err)
	}

	uc := usecase.NewPgUsecase(global.DB)
	restController := controller.NewHandler(uc)

	router := gin.Default()
	// router.Use(middleware.Cors())
	api.Bind(router, restController)

	port, _ := cmd.Flags().GetString("port")
	router.Run(":" + port)
}

func init() {
	Cmd.Flags().StringP("port", "p", "8080", "port")
	Cmd.Flags().
		StringVarP(&global.Database, `database`, "d", "postgres", `"postgres", "mysql"`)
	Cmd.Flags().
		StringVarP(&global.Env, `env`, "e", "dev", `"dev", "prod"`)

	Cmd.Flags().BoolVarP(&global.AutoMigrate, `automigrate`, "a", false, `"true", "false"`)
}
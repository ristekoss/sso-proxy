package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	// Blank-import the function package so the init() runs
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/ristekoss/ssoui-rest-proxy"
	"github.com/spf13/viper"
)

func main() {
	os.Setenv("FUNCTION_TARGET", "Proxy")

	err := LoadConfig()
	if err != nil {
		log.Fatal().AnErr("LoadConfig", err)
	}

	if viper.GetString("APP_ENV") == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		log.Logger = log.Output(os.Stdout)
	}

	if viper.GetBool("DISABLE_LOGGING") {
		log.Logger = log.Output(io.Discard)
	}

	zerolog.SetGlobalLevel(zerolog.Level(viper.GetInt("LOG_LEVEL")))

	if err := funcframework.Start(viper.GetString("PORT")); err != nil {
		log.Fatal().AnErr("funcframework.Start", err)
	}
}

func LoadConfig() (err error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.AllowEmptyEnv(true)

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("APP_ENV", "dev")
	viper.SetDefault("LOG_LEVEL", int(zerolog.InfoLevel))
	viper.SetDefault("DISABLE_LOGGING", false)
	viper.SetDefault("SERVICE_URL", "http://localhost:8081/")
	viper.SetDefault("CAS_URL", "https://sso.ui.ac.id/cas2/")

	err = viper.ReadInConfig()
	return err
}

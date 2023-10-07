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
	config, err := LoadConfig()
	if err != nil {
		log.Fatal().AnErr("LoadConfig", err)
	}

	if config.AppEnv == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		log.Logger = log.Output(os.Stdout)
	}

	if config.DisableLogging {
		log.Logger = log.Output(io.Discard)
	}

	zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))

	if err := funcframework.Start(config.Port); err != nil {
		log.Fatal().AnErr("funcframework.Start", err)
	}
}

type Config struct {
	Port           string `mapstructure:"PORT"`
	AppEnv         string `mapstructure:"APP_ENV"`
	LogLevel       int    `mapstructure:"LOG_LEVEL"`
	DisableLogging bool   `mapstructure:"DISABLE_LOGGING"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// default config
	config = Config{
		Port:           "8080",
		AppEnv:         "dev",
		LogLevel:       int(zerolog.InfoLevel),
		DisableLogging: false,
	}

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}

package config

import (
	"bytes"
	_ "embed"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

//go:embed default.env
var defaultConfig []byte

type Config struct {
	PSQL       DatabaseConfig `json:"psql" mapstructure:"psql" validate:"required"`
}

type DatabaseConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
	Database string `json:"database" mapstructure:"database"`
}

func (o Config) Validate() error {
	validate := validator.New()
	err := validate.Struct(o)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return err
		}
	}
	return nil
}

func LoadConfig() *Config {
	cfg := &Config{}

	v := viper.NewWithOptions(viper.KeyDelimiter("__"))

	v.AutomaticEnv()
	v.AddConfigPath(".")
	v.SetConfigType("env")
	v.SetConfigFile(".env")
	v.AddConfigPath("/app/config")
	v.AddConfigPath(".")
	v.SetConfigName("config")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	err := v.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		panic(err)
	}

	if err := v.ReadInConfig(); err != nil {
		log.Println("failed to read config from file")
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	if err = cfg.Validate(); err != nil {
		panic(err)
	}

	return cfg
}

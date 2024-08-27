package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Envs *Environments

// Environments define the environment variables
type Environments struct {
	APIPort             string `mapstructure:"API_PORT" envconfig:"API_PORT"`
	AppName             string `mapstructure:"APP_NAME" envconfig:"APP_NAME"`
	AppEnv              string `mapstructure:"APP_ENV" envconfig:"APP_ENV"`
	DatabaseName        string `mapstructure:"DATABASE_NAME" envconfig:"DATABASE_NAME"`
	DatabaseURL         string `mapstructure:"DATABASE_URL" envconfig:"DATABASE_URL"`
	DatabaseTimeout     int16  `mapstructure:"DATABASE_TIMEOUT" envconfig:"DATABASE_TIMEOUT"`
}

// LoadEnvVars load the environment variables
func LoadEnvVars() (*Environments, error) {
	godotenv.Load()
	c := &Environments{}
	if err := envconfig.Process("", c); err != nil {
		return nil, err
	}
	return c, nil
}

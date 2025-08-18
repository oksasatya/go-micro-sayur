package config

import "github.com/spf13/viper"

type App struct {
	AppPort      string `json:"app_port"`
	AppEnv       string `json:"app_env"`
	JWTSecretKey string `json:"jwt_secret_key"`
	JWTIssuer    string `json:"jwt_issuer"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DbName    string `json:"db_name"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type Config struct {
	App    App    `json:"app"`
	PsqlDB PsqlDB `json:"psql_db"`
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort:      viper.GetString("APP_PORT"),
			AppEnv:       viper.GetString("APP_ENV"),
			JWTSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JWTIssuer:    viper.GetString("JWT_ISSUER"),
		},
		PsqlDB: PsqlDB{
			Host:      viper.GetString("DATABASE_HOST"),
			Port:      viper.GetString("DATABASE_PORT"),
			User:      viper.GetString("DATABASE_USERNAME"),
			Password:  viper.GetString("DATABASE_PASSWORD"),
			DbName:    viper.GetString("DATABASE_NAME"),
			DBMaxOpen: viper.GetInt("DATABASE_MAX_OPEN_CONNECTIONS"),
			DBMaxIdle: viper.GetInt("DATABASE_MAX_IDLE_CONNECTIONS"),
		},
	}
}

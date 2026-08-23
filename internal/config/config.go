package config

import "github.com/spf13/viper"

type Database string

const (
	Postgres Database = "postgresql"
	Mysql    Database = "mysql"
	SQLite   Database = "sqlite"
)

type Profile struct {
	Name     string
	Server   string
	Port     string
	User     string
	Password string
}

func LoadConfig(path string) (prof Profile, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&prof)
	return
}

package config

import "github.com/spf13/viper"

type Config struct {
	General  General  `mapstructure:"general"`
	Database Database `mapstructure:"database"`
	Auth     Auth     `mapstructure:"auth"`
}

type General struct {
	AppEnv   string `mapstructure:"app_env"`
	HTTPPort string `mapstructure:"http_port"`
}

type Database struct {
	Postgresql Postgresql `mapstructure:"postgresql"`
	Migration  Migration  `mapstructure:"migration"`
}

type Postgresql struct {
	URI string `mapstructure:"uri"`
}

type Migration struct {
	Path string `mapstructure:"path"`
}

type Auth struct {
	JWTSecretKey string `mapstructure:"jwt_secret_key"`
}

func setDefault(v *viper.Viper) {
	v.SetDefault("general.app_env", DefaultAppEnv)
	v.SetDefault("general.http_port", DefaultHTTPPort)
	v.SetDefault("database.postgresql.uri", DefaultPostgresqlWalletURI)
	v.SetDefault("database.migration.path", DefaultMigrationsPath)
	v.SetDefault("auth.jwt_secret_key", DefaultJWTSecretKey)
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	setDefault(viper.GetViper())

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

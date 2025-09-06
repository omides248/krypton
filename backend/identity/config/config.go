package config

import "github.com/spf13/viper"

type Config struct {
	General      General      `mapstructure:"general"`
	MinIO        MinIO        `mapstructure:"minio"`
	Database     Database     `mapstructure:"database"`
	LocalStorage LocalStorage `mapstructure:"local_storage"`
	Auth         Auth         `mapstructure:"auth"`
}

type General struct {
	AppEnv   string `mapstructure:"app_env"`
	GRPCPort string `mapstructure:"grpc_port"`
	HTTPPort string `mapstructure:"http_port"`
	GRPCAddr string `mapstructure:"grpc_addr"`
	Host     string `mapstructure:"host"`
}

type MinIO struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	PublicURL string `mapstructure:"public_url"`
}

type Database struct {
	Postgresql Postgresql `mapstructure:"postgresql"`
	Migration  Migration  `mapstructure:"migration"`
}

type Postgresql struct {
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Migration struct {
	Path string `mapstructure:"path"`
}

type LocalStorage struct {
	PublicStoragePath  string `mapstructure:"public_storage_path"`
	PrivateStoragePath string `mapstructure:"private_storage_path"`
	StaticFilesPrefix  string `mapstructure:"static_files_prefix"`
}

type Auth struct {
	JWTSecretKey string `mapstructure:"jwt_secret_key"`
}

func setDefault(v *viper.Viper) {
	v.SetDefault("general.app_env", DefaultAppEnv)
	v.SetDefault("general.grpc_port", DefaultGRPCPort)
	v.SetDefault("general.http_port", DefaultHTTPPort)
	v.SetDefault("general.grpc_addr", DefaultGRPCAddr)
	v.SetDefault("general.host", DefaultHost)

	v.SetDefault("minio.endpoint", DefaultMinIOEndpoint)
	v.SetDefault("minio.access_key", DefaultMinIOAccessKey)
	v.SetDefault("minio.secret_key", DefaultMinIOSecretKey)
	v.SetDefault("minio.public_url", DefaultMinIOPublicURL)

	v.SetDefault("database.postgresql.uri", DefaultPostgresqlIdentityURI)
	v.SetDefault("database.migration.path", DefaultMigrationsPath)

	v.SetDefault("local_storage.public_storage_path", DefaultPublicStoragePath)
	v.SetDefault("local_storage.private_storage_path", DefaultPrivateStoragePath)
	v.SetDefault("local_storage.static_files_prefix", DefaultStaticFilesPrefix)

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

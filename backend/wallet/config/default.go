package config

var (
	DefaultAppEnv              = "development"
	DefaultHTTPPort            = ":8080"
	DefaultPostgresqlWalletURI = "postgres://omides248:123123@127.0.0.1:5432/wallet_db?sslmode=disable"
	DefaultMigrationsPath      = "file://internal/adapters/storage/postgresql/migrations"
	DefaultJWTSecretKey        = "25731f98a3959cc09469b86ffcff0e35702fe43337dd18f80cf9eae2767876f3"
)

package config

var (
	DefaultAppEnv   = "development"
	DefaultGRPCPort = ":50051"
	DefaultHTTPPort = ":8080"
	DefaultGRPCAddr = "127.0.0.1:50051"
	DefaultHost     = "192.168.8.140:8080"

	DefaultMinIOEndpoint  = "192.168.8.140:9000"
	DefaultMinIOAccessKey = "minioadmin"
	DefaultMinIOSecretKey = "minioadmin123123"
	DefaultMinIOPublicURL = "http://192.168.8.140:9000"

	DefaultPostgresqlIdentityURI = "postgres://omides248:123123@127.0.0.1:5432/identity_db"

	DefaultMigrationsPath = "file://internal/adapters/storage/postgresql/migration"

	DefaultPublicStoragePath  = "/var/lib/blackshop/storage/public"
	DefaultPrivateStoragePath = "/var/lib/blackshop/storage/private"
	DefaultStaticFilesPrefix  = "public"

	DefaultJWTSecretKey = "25731f98a3959cc09469b86ffcff0e35702fe43337dd18f80cf9eae2767876f3"
)

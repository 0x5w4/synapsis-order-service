package config

import (
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/spf13/viper"
)

type Config struct {
	App      *AppConfig
	Tracer   *TracerConfig
	HTTP     *HTTPConfig
	Postgres *DatabaseConfig
	GRPC     *GRPCConfig
}

type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Debug       bool
	UsePubsub   bool
	FrontendURL string
}

type TracerConfig struct {
	ServerURL      string
	SecretToken    string
	ServiceName    string
	ServiceVersion string
	Environment    string
	NodeName       string
}

type DatabaseConfig struct {
	DSN                string
	MigrateDSN         string
	DBName             string
	MaxOpenConns       int
	MaxIdleConns       int
	ConnMaxLifetime    int
	SlowQueryThreshold int
	Debug              bool
}

type HTTPConfig struct {
	Host               string
	Port               int
	BasePath           string
	DomainName         string
	EnableMigrationAPI bool
}

type GRPCConfig struct {
	InventoryHost string
	InventoryPort int
}

func LoadConfig(envPath string) (*Config, error) {
	if envPath == "" {
		envPath = ".env"
	}

	viper.SetConfigFile(envPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		var cfgErr viper.ConfigFileNotFoundError
		if !errors.As(err, &cfgErr) {
			return nil, err
		}
	}

	config := &Config{
		App: &AppConfig{
			Name:        viper.GetString("APP_NAME"),
			Version:     viper.GetString("APP_VERSION"),
			Environment: viper.GetString("APP_ENVIRONMENT"),
			Debug:       viper.GetBool("APP_DEBUG"),
			UsePubsub:   viper.GetBool("APP_USE_PUBSUB"),
			FrontendURL: viper.GetString("FRONTEND_URL"),
		},
		Tracer: &TracerConfig{
			ServerURL:      viper.GetString("ELASTIC_APM_SERVER_URL"),
			SecretToken:    viper.GetString("ELASTIC_APM_SECRET_TOKEN"),
			ServiceName:    viper.GetString("ELASTIC_APM_SERVICE_NAME"),
			ServiceVersion: viper.GetString("ELASTIC_APM_SERVICE_VERSION"),
			Environment:    viper.GetString("ELASTIC_APM_ENVIRONMENT"),
			NodeName:       viper.GetString("ELASTIC_APM_NODE_NAME"),
		},
		HTTP: &HTTPConfig{
			Host:               viper.GetString("HTTP_HOST"),
			Port:               viper.GetInt("HTTP_PORT"),
			BasePath:           viper.GetString("HTTP_BASE_PATH"),
			DomainName:         viper.GetString("HTTP_DOMAIN_NAME"),
			EnableMigrationAPI: viper.GetBool("HTTP_ENABLE_MIGRATION_API"),
		},
		Postgres: &DatabaseConfig{
			DSN:                viper.GetString("POSTGRES_DSN"),
			MigrateDSN:         viper.GetString("POSTGRES_MIGRATE_DSN"),
			DBName:             viper.GetString("POSTGRES_DB_NAME"),
			MaxOpenConns:       viper.GetInt("POSTGRES_MAX_OPEN_CONNS"),
			MaxIdleConns:       viper.GetInt("POSTGRES_MAX_IDLE_CONNS"),
			ConnMaxLifetime:    viper.GetInt("POSTGRES_CONN_MAX_LIFETIME"),
			SlowQueryThreshold: viper.GetInt("POSTGRES_SLOW_QUERY_THRESHOLD"),
			Debug:              viper.GetBool("POSTGRES_DEBUG"),
		},
		GRPC: &GRPCConfig{
			InventoryHost: viper.GetString("GRPC_INVENTORY_HOST"),
			InventoryPort: viper.GetInt("GRPC_INVENTORY_PORT"),
		},
	}

	return config, nil
}

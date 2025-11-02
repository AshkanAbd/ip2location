package config

import (
	"ip_location/cmd/http/config"
	pkgPgSql "ip_location/pkg/pgsql"
)

type AppConfig struct {
	HttpConfig  config.HTTPConfig `mapstructure:"http"`
	PgSQLConfig pkgPgSql.Config   `mapstructure:"pgsql"`
	LogLevel    string            `mapstructure:"log_level"`
}

package main

import (
	"ip_location/config"
	"ip_location/internal/http/handlers"
	"ip_location/internal/http/middlewares"
	"ip_location/internal/iptolocation/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pkgCfg "ip_location/pkg/config"
	pkgLog "ip_location/pkg/logger"
	pkgPgSql "ip_location/pkg/pgsql"
)

var Config config.AppConfig

func init() {
	pkgLog.SetLogLevel("Trace")
	Config = *pkgCfg.Load[config.AppConfig]()
	pkgLog.SetLogLevel(Config.LogLevel)
}

func main() {
	pgSqlConn, err := pkgPgSql.NewConnector(Config.PgSQLConfig)
	if err != nil {
		pkgLog.Error(err, "Failed to connect to pgSql")
		return
	}
	defer pgSqlConn.Close()

	if err := pgSqlConn.Migrate(); err != nil {
		pkgLog.Error(err, "Failed to migrate database")
		return
	}

	cfg := &gorm.Config{}
	gormConn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: pgSqlConn.GetConnection(),
	}), cfg)
	if err != nil {
		pkgLog.Error(err, "Failed to initialize gorm connection")
		return
	}

	ipToLocationSrv := services.NewIPToLocation(&http.Client{}, gormConn)

	httpHandler := handlers.NewHttpHandler(ipToLocationSrv)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middlewares.Logger())

	engine.GET("/metrics", handlers.Metrics())
	engine.GET("/health", handlers.HealthCheck())

	httpHandler.RegisterRoutes(engine)

	if err := engine.Run(Config.HttpConfig.Address); err != nil {
		pkgLog.Error(err, "Failed to start http server")
		return
	}
}

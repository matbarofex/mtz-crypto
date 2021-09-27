package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matbarofex/mtz-crypto/pkg/config"
	"github.com/matbarofex/mtz-crypto/pkg/controller"
	"github.com/matbarofex/mtz-crypto/pkg/crypto/cryptonator"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/service"
	"github.com/matbarofex/mtz-crypto/pkg/store/db"
	"github.com/matbarofex/mtz-crypto/pkg/store/memory"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.NewConfig(fs)

	// Gin mode
	if !cfg.GetBool("crypto.debug.mode") {
		gin.SetMode(gin.ReleaseMode)
	}

	// Start Gin Engine
	r := gin.New()

	// Conexión a DB
	gormDB := createGomDB(cfg)
	defer closeGormDBConnection(gormDB)

	// Stores
	walletStore := db.NewWalletStore(gormDB)
	marketDataStore := memory.NewMarketDataStore()

	// Market Data channel
	mdChannel := make(model.MdChannel)

	// Services
	marketDataService := service.NewMarketDataService(marketDataStore)
	walletService := service.NewWalletService(walletStore, marketDataService)

	// Start MD consumption
	marketDataService.ConsumeMD(mdChannel)

	// Cliente API externa
	httpClient := &http.Client{}
	cryptoClient := cryptonator.NewCryptonatorClient(cfg, httpClient, mdChannel)
	go cryptoClient.Start()

	// Controllers
	walletController := controller.NewWalletController(walletService)

	// Controller routes
	r.GET("/wallet/value", walletController.GetWalletValue)

	// Health check handler
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	// Load server configuration
	addr := cfg.GetString("crypto.http.addr")

	// HTTP Server
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start server
	go func() {
		// TODO log
		fmt.Println("starting HTTP server", addr)

		// TODO tratamiento del error
		_ = srv.ListenAndServe()
	}()

	// Wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// TODO log de la señal
	<-quit

	// Shutdown
	timeout := cfg.GetDuration("crypto.http.shutdown.timeout")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		// TODO log de error
		panic(err)
	}
}

// createGomDB configuración de acceso a datos y GORM
func createGomDB(cfg *config.Config) *gorm.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s %s",
		cfg.GetString("crypto.postgres.host"),
		cfg.GetInt("crypto.postgres.port"),
		cfg.GetString("crypto.postgres.username"),
		cfg.GetString("crypto.postgres.dbname"),
		cfg.GetString("crypto.postgres.password"),
		cfg.GetString("crypto.postgres.extraopts"),
	)

	dbPool := &sql.DB{}
	dbPool.SetMaxIdleConns(cfg.GetInt("crypto.postgres.maxidleconns"))

	gormConfig := &gorm.Config{
		Logger:      logger.Discard,
		PrepareStmt: true,
		ConnPool:    dbPool,
	}

	if cfg.GetBool("crypto.postgres.sqldebug") {
		newLogger := logger.New(
			log.New(os.Stderr, "crypto", 0), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      false,       // Disable color
			},
		)
		gormConfig.Logger = newLogger
	}

	db, err := gorm.Open(postgres.Open(connStr), gormConfig)
	if err != nil {
		log.Fatalf("error trying to connect to DB: %v", err)
	}

	return db
}

// closeGormDBConnection cierra conexiones a DB relacional
func closeGormDBConnection(db *gorm.DB) {
	stmtManger, ok := db.ConnPool.(*gorm.PreparedStmtDB)

	if ok {
		for _, stmt := range stmtManger.Stmts {
			stmt.Close() // close the prepared statement
		}
	}

	dbLocal, err := db.DB()
	if err == nil {
		dbLocal.Close() //CloseDB
	}
}

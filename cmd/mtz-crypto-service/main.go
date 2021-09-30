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

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/matbarofex/mtz-crypto/pkg"
	"github.com/matbarofex/mtz-crypto/pkg/config"
	"github.com/matbarofex/mtz-crypto/pkg/controller"
	"github.com/matbarofex/mtz-crypto/pkg/crypto/cryptonator"
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/service"
	"github.com/matbarofex/mtz-crypto/pkg/store/db"
	"github.com/matbarofex/mtz-crypto/pkg/store/memory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.NewConfig(fs)

	// Zap Logger
	logger := createZapLogger(cfg)
	defer logger.Sync()

	logger.Info("starting service")

	// Gin mode
	if !cfg.GetBool("crypto.debug.mode") {
		gin.SetMode(gin.ReleaseMode)
	}

	// Start Gin Engine
	r := gin.New()

	if cfg.GetBool("crypto.debugmode") {
		r.Use(ginzap.Ginzap(logger, time.RFC3339Nano, true))
	}
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Conexión a DB
	gormDB := createGomDB(cfg)
	defer closeGormDBConnection(gormDB)

	// Stores
	walletStore := db.NewWalletStore(gormDB)
	marketDataStore := memory.NewMarketDataStore()

	// Market Data channel
	mdChannel := make(model.MdChannel)

	// Services
	marketDataService := service.NewMarketDataService(logger, marketDataStore)
	walletService := service.NewWalletService(walletStore, marketDataService)

	// Start MD consumption
	marketDataService.ConsumeMD(mdChannel)

	// Cliente API externa
	cryptonatorHTTPClient := &http.Client{Timeout: cfg.GetDuration("crypto.api.cryptonator.timeout")}
	cryptoClient := cryptonator.NewCryptonatorClient(cfg, logger, cryptonatorHTTPClient, mdChannel)
	go cryptoClient.Start()

	// Controllers
	walletController := controller.NewWalletController(logger, walletService)

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
		logger.Info("starting HTTP server", zap.String("addr", addr))

		err := srv.ListenAndServe()
		if err == http.ErrServerClosed {
			logger.Info("shutting down server", zap.Error(err))
		} else {
			logger.Fatal("shutting down server", zap.Error(err))
		}
	}()

	// Wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	logger.Info("signal received", zap.Any("signal", s))

	// Shutdown
	timeout := cfg.GetDuration("crypto.http.shutdown.timeout")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("fatal error", zap.Error(err))
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

// createZapLogger inicializa logger de la app
func createZapLogger(cfg *config.Config) *zap.Logger {
	initialFields := make(map[string]interface{})
	if cfg.GetString("crypto.logging.format") == "json" {
		initialFields["svc"] = pkg.ServiceName
		initialFields["vsn"] = pkg.Version
	}

	level := zapcore.InfoLevel
	if cfg.GetBool("crypto.debug.mode") {
		level = zapcore.DebugLevel
	}

	configLogger := zap.Config{
		Encoding:         cfg.GetString("crypto.logging.format"),
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "timestamp",
			EncodeTime:  zapcore.RFC3339NanoTimeEncoder,
		},
		InitialFields: initialFields,
	}

	logger, err := configLogger.Build()
	if err != nil {
		log.Fatalf("error building zap logger: %s", err.Error())
	}

	return logger
}

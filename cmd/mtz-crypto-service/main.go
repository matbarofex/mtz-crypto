package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/matbarofex/mtz-crypto/pkg/config"
	"github.com/matbarofex/mtz-crypto/pkg/crypto/cryptonator"
	"github.com/matbarofex/mtz-crypto/pkg/store/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.NewConfig(fs)

	// Conexión a DB
	gormDB := createGomDB(cfg)
	defer closeGormDBConnection(gormDB)

	// Stores
	walletStore := db.NewWalletStore(gormDB)

	// TODO - eliminar
	fmt.Println("-----------------------------------------")
	wallet1, err := walletStore.GetWallet("wallet1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Wallet 1: %+v\n", wallet1)
	fmt.Println("-----------------------------------------")

	// Cliente API externa
	httpClient := &http.Client{}
	cryptoClient := cryptonator.NewCryptonatorClient(cfg, httpClient)
	cryptoClient.Start()
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

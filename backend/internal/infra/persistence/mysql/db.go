package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Gorm  *gorm.DB
	SQLDB *sql.DB
}

func NewMySQL(ctx context.Context, dsn string) (*DB, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql with gorm: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db from gorm: %w", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping mysql with gorm: %w", err)
	}

	if err := autoMigrate(gormDB); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("auto migrate with gorm: %w", err)
	}

	return &DB{
		Gorm:  gormDB,
		SQLDB: sqlDB,
	}, nil
}

func autoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.AppMetadata{},
		&model.Auction{},
		&model.Bid{},
		&model.SyncCursor{},
		&model.ProcessedLog{},
	); err != nil {
		return err
	}

	meta := model.AppMetadata{
		MetaKey: "schema_version",
		Value:   "001",
	}

	return db.Where(&model.AppMetadata{
		MetaKey: "schema_version",
	}).FirstOrCreate(&meta).Error
}

func (db *DB) Close() error {
	if db == nil || db.SQLDB == nil {
		return nil
	}
	return db.SQLDB.Close()
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"messenger/config"
	"messenger/internal/loggerzap"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection(cfg *config.ConfigDB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.GetUser(), cfg.GetPassword(), cfg.GetHost(),
		cfg.GetPort(), cfg.GetDBName(), cfg.GetSSLMode())

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	configurePool(sqlDB)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		return nil, err
	}

	loggerZap, err := loggerzap.NewLogger()
	if err != nil {
		sqlDB.Close()
		return nil, err
	}

	zapGORM := zapgorm2.New(loggerZap.GetLogger())

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		WithoutQuotingCheck:  true,
		PreferSimpleProtocol: true,
		Conn:                 sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: false,
		Logger:      zapGORM,
	})

	if err != nil {
		sqlDB.Close()
		return nil, err
	}

	return gormDB, nil
}

func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(1 * time.Hour)
	db.SetConnMaxIdleTime(5 * time.Minute)
}

package gorm

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	PGHost     string `long:"pg-host" env:"PG_HOST" required:"true"`
	PGUser     string `long:"pg-user" env:"PG_USER" required:"true"`
	PGPassword string `long:"pg-password" env:"PG_PASSWORD" required:"true"`
	PGName     string `long:"pg-name" env:"PG_NAME" required:"true"`
	PGPort     string `long:"pg-port" env:"PG_PORT" required:"true"`
}

func New(cfg Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.PGHost, cfg.PGUser, cfg.PGPassword, cfg.PGName, cfg.PGPort)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		),
	})
}

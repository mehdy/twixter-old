package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Twitter struct {
	db     *gorm.DB
	logger entities.Logger
}

func NewTwitter(config entities.ConfigGetter, logger entities.Logger) *Twitter {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		config.GetString("POSTGRES_USER"),
		config.GetString("POSTGRES_PASSWORD"),
		config.GetString("POSTGRES_HOST"),
		config.GetInt("POSTGRES_PORT"),
		config.GetString("POSTGRES_DB"),
		config.GetString("POSTGRES_SSLMODE"),
	)

	logger.As("I").Logf(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.As("F").WithError(err).Logf("Failed to connect to postgres db")
	}

	return &Twitter{
		db:     db,
		logger: logger,
	}
}

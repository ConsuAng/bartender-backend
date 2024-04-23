package database

import (
	"fmt"

	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Module = fx.Provide(NewDB)

func NewDB() (*gorm.DB, error) {
	fmt.Println("Conectando a la base de datos usando GORM...")
	dsn := "user:userpassword@tcp(localhost:3307)/angelusbartender?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object from GORM: %w", err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

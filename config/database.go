package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (c *Config) ConnectDatabase() (*gorm.DB, error) {
	// MySQL connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	log.Printf("Connecting to MySQL database: %s@%s:%s/%s",
		c.Database.Username, c.Database.Host, c.Database.Port, c.Database.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	log.Println("MySQL database connected successfully")
	return db, nil
} 
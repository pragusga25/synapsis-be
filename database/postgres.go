package database

import (
	"fmt"
	"log"
	"synapsis/config"
	"synapsis/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.Cart{}, &models.Order{}, &models.ProductCategory{}, &models.OrderItem{}, &models.Transaction{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Drop foreign keys
	db.Exec("ALTER TABLE product_categories DROP CONSTRAINT fk_product_categories_product")
	db.Exec("ALTER TABLE product_categories DROP CONSTRAINT fk_product_categories_category")

	// Add foreign keys with ON DELETE CASCADE
	db.Exec("ALTER TABLE product_categories ADD CONSTRAINT fk_product_categories_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE")
	db.Exec("ALTER TABLE product_categories ADD CONSTRAINT fk_product_categories_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE")
}

func InitDB(cfg config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Run Migrations
	AutoMigrate(DB)
}

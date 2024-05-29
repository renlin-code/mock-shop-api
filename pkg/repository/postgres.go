package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable           = "users"
	ordersTable          = "orders"
	productsTable        = "products"
	orderedProductsTable = "ordered_products"
	categoriesTables     = "categories"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)
	fmt.Println(conn)
	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

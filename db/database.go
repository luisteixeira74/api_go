package db

import (
	"apis/configs"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	db, err := sql.Open("sqlite3", config.DBFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao testar conexão com o banco: %w", err)
	}

	return db, nil
}

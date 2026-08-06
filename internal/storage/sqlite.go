package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar no banco: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("erro ao criar tabelas: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pessoa (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS cartao (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		tipo INTEGER NOT NULL, -- 0 para crédito, 1 para débito
		dia_fechamentO INTEGER,
		dia_vencimento INTEGER 
	);

	CREATE TABLE IF NOT EXISTS gasto (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		descricao TEXT NOT NULL,
		valor_total REAL NOT NULL,
		data_compra TEXT NOT NULL, -- Formato YYYY-MM-DD
		total_parcelas INTEGER NOT NULL,
		pessoa_id INTEGER,
		cartao_id INTEGER,
		FOREIGN KEY(pessoa_id) REFERENCES pessoa(id),
		FOREIGN KEY(cartao_id) REFERENCES cartao(id)
	);

	CREATE TABLE IF NOT EXISTS parcela (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		gasto_id INTEGER,
		numero_parecela INTEGER NOT NULL,
		valor REAL NOT NULL,
		data_cobranca TEXT NOT NULL, -- Formato YYYY-MM-DD
		status_pagamento INTEGER NOT NULL, -- 0 para pendente, 1 para pago
		FOREIGN KEY(gasto_id) REFERENCES gasto(id)
	);
	`

	_, err := db.Exec(query)
	return err
}

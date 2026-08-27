package repository

import (
	"fmt"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB inicializa a conexão com o banco de dados, criando-o em dbPath.
// Retorna um ponteiro para o bd criado e erros de abertura, conexão ou criaçã do banco de dados.
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

// createTables cria todas as tabelas de entidades definidas na modelagem
// Retorna erro em caso de falha de execução.
func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS persons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		active INTEGER NOT NULL -- 0 para inativo, 1 para ativo
	);

	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		active INTEGER NOT NULL, -- 0 para inativo, 1 para ativo
		type INTEGER NOT NULL, -- 0 para crédito, 1 para débito
		closing_day INTEGER,
		due_day INTEGER 
	);

	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		total_value REAL NOT NULL,
		purchase_date TEXT NOT NULL, -- Formato YYYY-MM-DD
		total_installments INTEGER NOT NULL,
		person_id INTEGER,
		card_id INTEGER,
		FOREIGN KEY(person_id) REFERENCES persons(id),
		FOREIGN KEY(card_id) REFERENCES cards(id)
	);

	CREATE TABLE IF NOT EXISTS installments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expense_id INTEGER,
		number_installments INTEGER NOT NULL,
		value REAL NOT NULL,
		due_date TEXT NOT NULL, -- Formato YYYY-MM-DD
		payment_status INTEGER NOT NULL, -- 0 para pendente, 1 para pago
		FOREIGN KEY(expense_id) REFERENCES expenses(id)
	);
	`

	_, err := db.Exec(query)
	return err
}

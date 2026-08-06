package storage

import (
	"database/sql"
	"gfinancer/internal/domain"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erro ao abrir banco em memória: %v", err)
	}

	query := `
		CREATE TABLE expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			description TEXT,
			total_value REAL,
			total_installments INTEGER,
			purchase_date TEXT,
			person_id INTEGER,
			card_id INTEGER
		);
		CREATE TABLE installments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			expense_id INTEGER,
			number_installment INTEGER,
			value REAL,
			due_date TEXT,
			payment_status INTEGER
		);
	`
	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("Erro ao criar tabelas no banco em memória: %v", err)
	}

	return db
}

func TestExpenseRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewExpense(db)

	t.Run("Sucesso: Insere despesa e 3 parcelas", func(t *testing.T) {
		pessoaId := uint(1)

		despesa := &domain.Expense{
			Description:       "Compra de Teste",
			TotalValue:        100.00,
			TotalInstallments: 3,
			PurchaseDate:      "2026-08-01",
			PersonId:          &pessoaId,
		}

		err := repo.Save(despesa, nil)
		if err != nil {
			t.Fatalf("Não esperava erro ao salvar: %v", err)
		}

		var countDespesas int
		err = db.QueryRow("SELECT COUNT(*) FROM expenses WHERE id = ?", despesa.Id).Scan(&countDespesas)
		if err != nil || countDespesas != 1 {
			t.Errorf("A despesa não foi inserida no banco corretamente")
		}

		var countParcelas int
		err = db.QueryRow("SELECT COUNT(*) FROM installments WHERE expense_id = ?", despesa.Id).Scan(&countParcelas)
		if err != nil || countParcelas != 3 {
			t.Errorf("Esperava 3 parcelas inseridas, encontrou %d", countParcelas)
		}
	})

	t.Run("Rollback: Falha na geração das parcelas cancela a despesa", func(t *testing.T) {
		pessoaId := uint(1)

		despesaInvalida := &domain.Expense{
			Description:       "Compra Falha",
			TotalValue:        50.00,
			TotalInstallments: 2,
			PurchaseDate:      "01/08/2026",
			PersonId:          &pessoaId,
		}

		err := repo.Save(despesaInvalida, nil)
		if err == nil {
			t.Fatalf("Esperava um erro por causa da data inválida, mas nada ocorreu")
		}

		var countDespesasFalhas int
		db.QueryRow("SELECT COUNT(*) FROM expenses WHERE description = 'Compra Falha'").Scan(&countDespesasFalhas)

		if countDespesasFalhas > 0 {
			t.Errorf("O ROLLBACK FALHOU! A despesa 'Compra Falha' foi encontrada no banco.")
		}
	})
}

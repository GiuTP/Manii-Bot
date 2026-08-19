package repository

import (
	"database/sql"
	"gfinancer/internal/domain"
)

// ExpenseRepo é a estrutura que guarda o ponteiro da conexão com o banco de dados da tabela "expense"
type ExpenseRepo struct {
	db *sql.DB
}

// NewExpenseRepo é a função construtora de ExpenseRepo
func NewExpenseRepo(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

// Save recebe um gasto (expense) e o cartão que este gasto está atrelado e em seguida
// salva no banco de dados o gasto e o parcelamento.
// Retorna error se não foi possível concluir.
func (r *ExpenseRepo) Save(expense *domain.Expense, card *domain.Card) error {
	// Inicia a conexão com o banco de dados, abrindo uma transação
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Agendamento de rollback. Em caso de qualquer erro durante a execução, o rollback será executado.
	defer tx.Rollback()

	// Preparação e execução da query SQL em "expenses"
	queryExpense := `
		INSERT INTO expenses (description, total_value, total_installments, purchase_date, person_id, card_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := tx.Exec(queryExpense,
		expense.Description,
		expense.TotalValue,
		expense.TotalInstallments,
		expense.PurchaseDate,
		expense.PersonId,
		expense.CardId,
	)
	if err != nil {
		return err
	}

	// Pega o último gasto inserido (inserido logo acima) em "expenses"
	expenseId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	expense.Id = uint(expenseId)

	// Gerador de slice de parcelas
	installments, err := expense.InstallmentGenerate(card)
	if err != nil {
		return err
	}

	// Preparação de query SQL em "installments"
	queryInstallment := `
		INSERT INTO installments (expense_id, number_installment, value, due_date, payment_status)
		VALUES (?, ?, ?, ?, ?)
	`
	// Inserção de query SQL em "installments"
	for _, installment := range installments {
		_, err = tx.Exec(queryInstallment,
			installment.ExpenseId,
			installment.NumberInstallments,
			installment.Value,
			installment.DueDate,
			installment.PaymentStatus,
		)
		if err != nil {
			return err
		}
	}

	// Commit no banco de dados, encerrando a transação
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ExpenseRepo) GetAll() ([]domain.Expense, error) {
	query := `
		SELECT id, description, total_value, total_installments, purchase_data, person_id, card_id
		FROM expenses
		ORDER BY purchase_Date DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var expenses []domain.Expense

	for rows.Next() {
		var e domain.Expense

		err := rows.Scan(
			&e.Id,
			&e.Description,
			&e.TotalValue,
			&e.TotalInstallments,
			&e.PurchaseDate,
			&e.PersonId,
			&e.CardId,
		)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

// Update atualiza apenas **metadados**
func (r *ExpenseRepo) Update(expense *domain.Expense) error {
	query := `
		UPDATE expenses
		SET description = ?, person_id = ?
		WHERE id = ?
	`

	res, err := r.db.Exec(
		query,
		expense.Description,
		expense.PersonId,
		expense.Id,
	)
	if err != nil {
		return err
	}

	checkAffectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if checkAffectedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete apaga (hard) uma expense e as installments vinculadas à ela.
func (r *ExpenseRepo) Delete(id uint) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	queryDeleteInstallments := `DELELE FROM installments WHERE expense_id = ?`
	_, err = tx.Exec(queryDeleteInstallments, id)
	if err != nil {
		return err
	}

	queryDeleteExpense := `DELETE FROM expenses WHERE id = ?`
	res, err := tx.Exec(queryDeleteExpense, id)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return sql.ErrNoRows
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

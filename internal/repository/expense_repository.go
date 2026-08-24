package repository

import (
	"database/sql"

	"gfinancer/internal/domain"
)

// *****************************
// ** CRUD de repositório expense
// *****************************

// ExpenseRepo gerencia a persistência e operações de banco de dados para a entidade de despesas e suas parcelas.
type ExpenseRepo struct {
	db *sql.DB
}

// NewExpenseRepo cria e retorna uma nova instância de ExpenseRepo utilizando a conexão de banco de dados fornecida.
func NewExpenseRepo(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

// Save cria uma nova despesa e registra todas as suas respectivas parcelas em uma transação atômica.
// O ID autoincrementado gerado pelo banco é atribuído à própria estrutura da despesa fornecida.
// Retorna um erro caso a transação falhe ou ocorra problema na geração das parcelas.
func (r *ExpenseRepo) Save(expense *domain.Expense, card *domain.Card) error {
	// Inicia a conexão com o banco de dados, abrindo uma transação
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Agendamento de rollback, garantindo atomocidade.
	defer tx.Rollback()

	query := `
		INSERT INTO expenses (description, total_value, total_installments, purchase_date, person_id, card_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := tx.Exec(query,
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

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	expense.Id = uint(id)

	// Gerador de slice de parcelas
	installments, err := expense.InstallmentGenerate(card)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO installments (expense_id, number_installment, value, due_date, payment_status)
		VALUES (?, ?, ?, ?, ?)
	`
	// Inserção de query SQL em "installments"
	for _, installment := range installments {
		_, err = tx.Exec(query,
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

// GetAll busca todas as despesas cadastradas no banco de dados ordenadas pela data de compra mais recente para a mais antiga.
// Retorna um sliced com todas as despesas encontradas, ou um erro se a consulta falhar.
func (r *ExpenseRepo) GetAll() ([]domain.Expense, error) {
	query := `
		SELECT id, description, total_value, total_installments, purchase_date, person_id, card_id
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

// Update atualiza os metadados de uma despesa existente (descrição e pessoa associada) com base no seu ID.
// Retorna sql.ErrNoRows se despesa não existir, ou outro erro em caso de falha na execução.
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

// Delete remove permanentemente (hard delete) uma despesa e todas as suas parcelas vinculadas em uma transação atômica.
// Retorna sql.ErrNoRows se a despesa informada não existir, ou um erro caso ocorra falha durante a transação.
func (r *ExpenseRepo) Delete(id uint) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Apaga de installments
	queryDeleteInstallments := `DELELE FROM installments WHERE expense_id = ?`
	_, err = tx.Exec(queryDeleteInstallments, id)
	if err != nil {
		return err
	}

	// Apaga de expenses
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

	// Commit
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

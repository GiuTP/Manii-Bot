package storage

import (
	"database/sql"
	"gfinancer/internal/domain"
)

type ExpenseRepo struct {
	db *sql.DB
}

func NewExpense(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Save(expense *domain.Expense, card *domain.Card) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

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

	expenseId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	expense.Id = uint(expenseId)

	installments, err := expense.InstallmentGenerate(card)
	if err != nil {
		return err
	}

	queryInstallment := `
		INSERT INTO installments (expense_id, number_installment, value, due_date, payment_status)
		VALUES (?, ?, ?, ?, ?)
	`

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

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

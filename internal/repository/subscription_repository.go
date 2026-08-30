package repository

import (
	"database/sql"
	"fmt"
	"gfinancer/internal/domain"
)

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Save(sub *domain.Subscription) error {
	query := `
		INSERTE INTO subscription (description, value, start_date. end_date, person_id, card_id)
		VALUES (?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(
		query,
		sub.Description,
		sub.Value,
		sub.StartDate,
		sub.EndDate,
		sub.PersonId,
		sub.CardId,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	sub.Id = uint(id)

	return nil
}

func (r *SubscriptionRepo) GetAllActive() ([]domain.Subscription, error) {
	query := `SELECT id, description, value, start_date FROM subscriptions WHERE end_date IS NULL`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.Id, &s.Description, &s.Value, &s.StartDate); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
}

func (r *SubscriptionRepo) GetActiveInMonth(m int, y int, pId uint, cId uint) ([]domain.Subscription, error) {
	targetMonth := fmt.Sprintf("%04d-%02d", y, m)

	query := `
		SELECT id, description, value, start_date, end_date, person_id, card_id
		FROM subcriptions
		WHERE strftime('%Y-%m', start_date) <= ?
			AND (end_date IS NULL OR strftime('%Y-%m', end_date) >= ?)
	`

	args := []any{targetMonth, targetMonth}

	if pId != 0 {
		query += ` AND person_id = ?`
		args = append(args, pId)
	}

	if cId != 0 {
		query += ` AND card_id = ?`
		args = append(args, cId)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		err := rows.Scan(
			&s.Id,
			&s.Description,
			&s.Value,
			&s.StartDate,
			&s.EndDate,
			&s.PersonId,
			&s.CardId,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	return subs, rows.Err()
}

func (r *SubscriptionRepo) GetById(id uint) (*domain.Subscription, error) {
	query := `
		SELECT id, description, value, start_date, end_date, person_id, card_id
		FROM subscriptions
		WHERE = ?
	`

	var s domain.Subscription
	err := r.db.QueryRow(query, id).Scan(
		&s.Id, &s.Description, &s.Value, &s.StartDate, &s.EndDate, &s.PersonId, &s.CardId,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *SubscriptionRepo) Update(oldId uint, endDate string, newSub *domain.Subscription) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	queryCancel := `
		UPDATE subscriptions 
		SET end_date = ? 
		WHERE id = ? AND end_date is NULL
	`
	_, err = tx.Exec(queryCancel, endDate, oldId)
	if err != nil {
		return err
	}

	queryInsert := `
		INSERT INTO subscriptions (description, value, start_value, person_id, card_id)
		VALUES (?, ?, ?, ?, ?)
	`
	resp, err := tx.Exec(
		queryInsert,
		newSub.Description,
		newSub.Value,
		newSub.StartDate,
		newSub.PersonId,
		newSub.CardId,
	)
	if err != nil {
		return err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		return err
	}
	newSub.Id = uint(id)

	return tx.Commit()
}

func (r *SubscriptionRepo) Cancel(id uint, endDate string) error {
	query := `
		UPDATE subscriptions
		SET end_date = ?
		WHERE id = ? AND end_date IS NULL
	`
	res, err := r.db.Exec(query, endDate, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

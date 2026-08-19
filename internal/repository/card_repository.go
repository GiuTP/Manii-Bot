package repository

import (
	"database/sql"
	"gfinancer/internal/domain"
	"strings"
)

type CardRepo struct {
	db *sql.DB
}

func NewCard(db *sql.DB) *CardRepo {
	return &CardRepo{db: db}
}

func (r *CardRepo) Create(card *domain.Card) error {
	query := `
		INSERT INTO cards(name, active, type, closing_day, due_day)
		VALUES (?, 1, ?, ?, ?)
	`

	res, err := r.db.Exec(
		query,
		card.Name,
		card.Type,
		card.ClosingDay,
		card.DueDay,
	)

	if err != nil {
		return err
	}

	cardId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	card.Id = uint(cardId)

	return nil
}

func (r *CardRepo) ReadMap() (map[string]domain.Card, error) {
	query := `SELECT id, name, closing_day, due_day, type FROM cards`

	rows, err := r.db.Query(query)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	cardsMap := make(map[string]domain.Card)

	for rows.Next() {
		var c domain.Card

		err := rows.Scan(&c.Id, &c.Name, &c.ClosingDay, &c.DueDay, &c.Type)
		if err != nil {
			return nil, err
		}

		chave := strings.ToLower(c.Name)
		cardsMap[chave] = c
	}

	return cardsMap, nil
}

func (r *CardRepo) Update(card *domain.Card) error {
	query := `
		UPDATE cards
		SET closing_day = ?, due_day = ?
		WHERE id = ?
	`

	res, err := r.db.Exec(
		query,
		card.ClosingDay,
		card.DueDay,
		card.Id,
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

func (r *CardRepo) Disable(id uint) error {
	query := `
		UPDATE cards SET active = ? WHERE id = ?
	`

	res, err := r.db.Exec(query, 0, id)
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

func (r *CardRepo) Enable(id uint) error {
	query := `
		UPDATE cards SET active = ? WHERE id = ?
	`

	res, err := r.db.Exec(query, 1, id)
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

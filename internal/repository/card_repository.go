package repository

import (
	"strings"

	"database/sql"

	"gfinancer/internal/domain"
)

// *****************************
// ** CRUD de repositório card
// *****************************

// CardRepo gerencia a persistência e operações de banco de dados para a entidade de cartões
type CardRepo struct {
	db *sql.DB
}

// NewCardRepo cria e retorna uma nova instância de CardRepo utilizando a conexão de banco de dados fornecida.
func NewCardRepo(db *sql.DB) *CardRepo {
	return &CardRepo{db: db}
}

// Create insere um novo cartão no banco de dados utilizando os dados da estrutura fornecida.
// A função define automaticamente o cartão como ativo (1) e atualiza a estrutura original com o ID gerado.
// Retorna um erro em caso de falha na inserção ou na recuperação do último ID inserido.
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

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	card.Id = uint(id)

	return nil
}

// ReadMap busca todos os cartões cadastrados e os retorna indexados em um mapa.
// A chave do mapa é o nome do cartão em letras minúsculas.
// Retorna o mapa populado em caso de sucesso, ou um erro se a consulta ao banco falhar.
func (r *CardRepo) ReadMap() (map[string]domain.Card, error) {
	query := `SELECT id, name, closing_day, due_day, type FROM cards WHERE active = 1`

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
		cardsMap[strings.ToLower(c.Name)] = c
	}

	return cardsMap, nil
}

// GetAll busca todos os cartões ativos cadastrados no banco de dados ordenados pelo nome.
// Retorna um sliced com todas os cartõtes encontradas, ou um erro se a consulta falhar.
func (r *CardRepo) GetAll() ([]domain.Card, error) {
	query := `
		SELECT id, name, active , type, closing_day, due_day 
		FROM cards
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cards []domain.Card
	for rows.Next() {
		var c domain.Card

		err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Active,
			&c.Type,
			&c.ClosingDay,
			&c.DueDay,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

// Update atualiza os dias de fechamento e vencimento de um cartão específico no banco de dados.
// Retorna um erro em caso de falha na execução, ou sql.ErrNoRows caso o cartão não exista.
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

// Disable inativa um cartão, alterando seu status para inativo (0) no banco de dados com base no ID.
// Retorna um erro em caso de falha na execução, ou sql.ErrNoRows caso o cartão não exista.
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

// Disable inativa um cartão previamente inativo, alterando seu status para ativo (1) no banco de dados com base no ID.
// Retorna um erro em caso de falha na execução, ou sql.ErrNoRows caso o cartão não exista.
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

package repository

import (
	"strings"

	"database/sql"

	"gfinancer/internal/domain"
)

// PersonRepo gerencia a persistência e operações de banco de dados para a entidade de pessoas.
type PersonRepo struct {
	db *sql.DB
}

// NewPersonRepo cria e retorna uma nova instância de PersonRepo utilizando a conexão de banco de dados fornecida.
func NewPersonRepo(db *sql.DB) *PersonRepo {
	return &PersonRepo{db: db}
}

// Create insere uma nova pessoa no banco de dados utilizando os dados da estrutura fornecida.
// A função define automaticamente a pessoa como ativa (1) e atualiza a estrutura original com o ID gerado.
// Retorna um erro em caso de falha na inserção ou na recuperação do último ID inserido.
func (r *PersonRepo) Create(person *domain.Person) error {
	query := `
		INSERT INTO persons(name, active)
		VALUES (?, 1)
	`

	res, err := r.db.Exec(
		query,
		person.Name,
	)
	if err != nil {
		return err
	}

	personId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	person.Id = uint(personId)

	return nil
}

// ReadMap busca todas as pessoas ativas cadastradas e as retorna indexadas em um mapa.
// A chave do mapa é o nome da pessoa em letras minúsculas.
// Retorna o mapa populado em caso de sucesso, ou um erro se a consulta ao banco falhar.
func (r *PersonRepo) ReadMap() (map[string]domain.Person, error) {
	query := `SELECT id, name FROM persons WHERE active = 1`

	rows, err := r.db.Query(query)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	personsMap := make(map[string]domain.Person)

	for rows.Next() {
		var p domain.Person

		err := rows.Scan(&p.Id, &p.Name)
		if err != nil {
			return nil, err
		}

		personsMap[strings.ToLower(p.Name)] = p
	}

	return personsMap, nil
}

// GetAll busca todas as pessoas ativas no banco de dados ordenadas alfabeticamente pelo nome.
// Retorna um slice com as pessoas encontradas, ou um erro se a consulta falhar.
func (r *PersonRepo) GetAll() ([]domain.Person, error) {
	query := `
		SELECT id, name, active
		FROM persons
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var persons []domain.Person
	for rows.Next() {
		var p domain.Person

		err = rows.Scan(
			&p.Id,
			&p.Name,
			&p.Active,
		)
		if err != nil {
			return nil, err
		}

		persons = append(persons, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

// Update atualiza o nome de uma pessoa existente no banco de dados com base no seu ID.
// Retorna sql.ErrNoRows caso nenhuma pessoa seja afetada (ID inexistente), ou outro erro em caso de falha na execução.
func (r *PersonRepo) Update(person *domain.Person) error {
	query := `
		UPDATE persons
		SET name = ?
		WHERE id = ?
	`

	res, err := r.db.Exec(
		query,
		person.Name,
		person.Id,
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

// Disable inativa uma pessoa, alterando seu status para inativo (0) no banco de dados com base no ID.
// Retorna sql.ErrNoRows caso a pessoa não seja encontrada, ou um erro se a execução falhar.
func (r *PersonRepo) Disable(id uint) error {
	query := `
		UPDATE persons SET active = ? WHERE id = ?
	`
	res, err := r.db.Exec(query, 0, id)
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

	return nil
}

// Enable reativa uma pessoa previamente inativada, restaurando seu status para ativo (1) no banco de dados.
// Retorna sql.ErrNoRows caso a pessoa não seja encontrada, ou um erro se a execução falhar.
func (r *PersonRepo) Enable(id uint) error {
	query := `
		UPDATE persons SET active = ? WHERE id = ?
	`
	res, err := r.db.Exec(query, 1, id)
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

	return nil
}

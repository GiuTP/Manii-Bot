package repository

import (
	"database/sql"
	"gfinancer/internal/domain"
	"strings"
)

// PersonRepo é a estrutura que guarda o ponteiro da conexão com o banco de dados da tabela "person"
type PersonRepo struct {
	db *sql.DB
}

// NewPersonRepo é a função construtora de PersonRepo
func NewPersonRepo(db *sql.DB) *PersonRepo {
	return &PersonRepo{db: db}
}

func (r *PersonRepo) Create(person *domain.Person) error {
	queryPerson := `
		INSERT INTO persons(name, active)
		VALUES (?, 1)
	`

	res, err := r.db.Exec(
		queryPerson,
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

// ReadMap carrega todos as pessoas salvas no banco de dados.
// Retorna um map de "person" e possíveis erros
func (r *PersonRepo) ReadMap() (map[string]domain.Person, error) {
	// Query SQL para pegar os id das pessoas inseridas no bd
	query := `SELECT id, name FROM persons`

	// Executa a query e devolve as linhas existentes em "person"
	rows, err := r.db.Query(query)
	if err != nil || rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	personsMap := make(map[string]domain.Person)

	// Carrega todos os nomes e ids de "person" no mapa de pessoas
	for rows.Next() {
		var p domain.Person

		err := rows.Scan(&p.Id, &p.Name)
		if err != nil {
			return nil, err
		}

		k := strings.ToLower(p.Name)
		personsMap[k] = p
	}

	return personsMap, nil
}

//func (r *PersonRepo) GetAll() ([]domain.Person, error){}

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

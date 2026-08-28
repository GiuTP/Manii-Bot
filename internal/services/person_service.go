package services

import (
	"errors"
	"fmt"
	"strings"

	"database/sql"

	"gfinancer/internal/domain"
	"gfinancer/internal/repository"
)

type PersonService struct {
	repo *repository.PersonRepo
}

func NewPersonService(r *repository.PersonRepo) *PersonService {
	return &PersonService{repo: r}
}

func (s *PersonService) Create(msg string) (string, error) {
	name := strings.TrimSpace(msg)
	if name == "" {
		return "", errors.New("Use: [nome]")
	}

	pMap, err := s.repo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao ler o mapa de pessoas: %w", err)
	}
	// Futuramente vou dá suporte a nome repetido usando chave composta.
	if _, exists := pMap[strings.ToLower(name)]; exists {
		return "", errors.New("Pessoa digitada já existe")
	}

	newP := &domain.Person{
		Name: name,
	}
	if err = s.repo.Create(newP); err != nil {
		return "", fmt.Errorf("Falha de criação no banco de dados: %w", err)
	}

	return fmt.Sprintf("%s cadastrado com sucesso!", newP.Name), nil
}

func (s *PersonService) List(msg string) (string, error) {
	flag := strings.TrimSpace(msg)
	if flag == "" {
		flag = "ativos"
	}

	persons, err := s.repo.GetAll()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar os registros de pessoa: %w", err)
	}
	if len(persons) == 0 {
		return "Nenhuma pessoa cadastrada.", nil
	}

	var sb strings.Builder
	if _, err = sb.WriteString(fmt.Sprintf("Pessoas %s: \n\n", flag)); err != nil {
		return "", fmt.Errorf("Falha de inclusão de texto: %w", err)
	}

	count := 0
	for _, p := range persons {
		if flag == "ativos" && p.Active == 0 {
			continue
		}

		if flag == "inativos" && p.Active == 1 {
			continue
		}

		status := ""
		if flag == "todos" && p.Active == 0 {
			status = "(Inativo)"
		}

		if _, err = sb.WriteString(fmt.Sprintf("#%d - %s %s\n", p.Id, p.Name, status)); err != nil {
			return "", fmt.Errorf("Falha de inclusão de texto: %w", err)
		}
		count++
	}

	if count == 0 {
		return fmt.Sprintf("Nenhuma pessoa encontrada com o filtro: %s", flag), nil
	}

	return sb.String(), nil
}

func (s *PersonService) Update(msg string) (string, error) {
	tokens := strings.SplitN(msg, " ", 2)
	if len(tokens) != 2 {
		return "", errors.New("Use: [nome_atual] [novo_nome]")
	}

	pMap, err := s.repo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa: %w", err)
	}

	person, exists := pMap[strings.ToLower(tokens[0])]
	if !exists {
		return "", errors.New("Pessoa não encontrada. Use /lista p ativos")
	}

	person.Name = tokens[1]

	if err = s.repo.Update(&person); err != nil {
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("Pessoa (ID %d) atualizada com sucesso!", person.Id), nil
}

func (s *PersonService) Disable(msg string) (string, error) {
	name := strings.Split(strings.TrimSpace(msg), " ")
	if len(name) != 1 {
		return "", errors.New("Use: [nome]")
	}

	pMap, err := s.repo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa: %w", err)
	}

	p, exists := pMap[strings.ToLower(name[0])]
	if !exists {
		return "", errors.New("Nome não encontrado. Use: /lista p inativos")
	}

	if err = s.repo.Disable(p.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Nome não encontrado ou já desativado. Use: /lista p inativos")
		}
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("%s (ID %d) desativado(a) com sucesso!", p.Name, p.Id), nil
}

func (s *PersonService) Enable(msg string) (string, error) {
	name := strings.Split(strings.TrimSpace(msg), " ")
	if len(name) != 1 {
		return "", errors.New("Use: [nome]")
	}

	pMap, err := s.repo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa: %w", err)
	}

	p, exists := pMap[strings.ToLower(name[0])]
	if !exists {
		return "", errors.New("Nome não encontrado. Use: /lista p ativos")
	}

	if err = s.repo.Enable(p.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Nome não encontrado ou já desativado. Use: /lista p ativos")
		}
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("%s (ID %d) reativado(a) com sucesso!", p.Name, p.Id), nil
}

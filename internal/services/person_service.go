package services

import (
	"errors"
	"fmt"
	"strconv"
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
		return "", errors.New("Use: [id] [novo_nome]")
	}

	id, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("Erro de conversão: %w", err)
	}

	pUpd := &domain.Person{
		Id:   uint(id),
		Name: tokens[1],
	}

	if err = s.repo.Update(pUpd); err != nil {
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("Pessoa (ID %d) atualizada com sucesso!", pUpd.Id), nil
}

func (s *PersonService) Disable(msg string) (string, error) {
	idT := strings.Split(strings.TrimSpace(msg), " ")
	if len(idT) != 1 {
		return "", errors.New("Use: [id]")
	}

	id, err := strconv.ParseUint(idT[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("Erro de conversão: %w", err)
	}

	if err = s.repo.Disable(uint(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("ID não encontrado ou já desativado. Utilize: /lista p inativos")
		}
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("Pessoa (ID %d) desativado(a) com sucesso!", id), nil
}

func (s *PersonService) Enable(msg string) (string, error) {
	idT := strings.Split(strings.TrimSpace(msg), " ")
	if len(idT) != 1 {
		return "", errors.New("Use: [id]")
	}

	id, err := strconv.ParseUint(idT[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("Erro de conversão: %w", err)
	}

	if err = s.repo.Enable(uint(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("ID não encontrado ou já ativo. Utilize: /lista p ativos")
		}
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("Pessoa (ID %d) ativado(a) com sucesso!", id), nil
}

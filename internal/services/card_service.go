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

type CardService struct {
	repo *repository.CardRepo
}

func NewCardService(r *repository.CardRepo) *CardService {
	return &CardService{repo: r}
}

func (s *CardService) Create(msg string) (string, error) {
	tokens := strings.Split(strings.TrimSpace(msg), " ")
	if len(tokens) != 4 {
		return "", errors.New("Use: [nome] [tipo] [dia_fechamento] [dia_vencimento]")
	}

	cMap, err := s.repo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa: %w", err)
	}

	// Separação de atributos
	name := strings.ToLower(tokens[0])
	// Futuramente vou dá suporte a nome repetido usando chave composta.
	if _, exists := cMap[name]; exists {
		return "", errors.New("Nome de cartão já usado.")
	}
	var typ uint8
	switch tokens[1] {
	case "credito", "crédito":
		typ = 0
	case "débito", "debito":
		typ = 1
	default:
		return "", errors.New("Tipo de cartão inválido. Tente: crédito ou débito")
	}
	closingD, err := strconv.ParseUint(tokens[2], 10, 8)
	if err != nil {
		return "", fmt.Errorf("Erro de conversão de fechamento: %w", err)
	}
	dueD, err := strconv.ParseUint(tokens[3], 10, 8)
	if err != nil {
		return "", fmt.Errorf("Erro de conversão de vencimento: %w", err)
	}
	newC := &domain.Card{
		Name:       name,
		Type:       typ,
		ClosingDay: uint8(closingD),
		DueDay:     uint8(dueD),
	}

	if err = s.repo.Create(newC); err != nil {
		return "", fmt.Errorf("Falha ao criar no banco de dados: %w", err)
	}

	return fmt.Sprintf("Cartão %s cadastrado com sucesso!", newC.Name), nil
}

func (s *CardService) List(msg string) (string, error) {
	flag := strings.TrimSpace(msg)
	if flag == "" {
		flag = "ativos"
	}

	cards, err := s.repo.GetAll()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar os registros de cartões: %w", err)
	}
	if len(cards) == 0 {
		return "Nenhum cartão cadastrado.", nil
	}

	var sb strings.Builder
	if _, err = sb.WriteString(fmt.Sprintf("Cartões %s: \n\n", flag)); err != nil {
		return "", fmt.Errorf("Falha de inclusão de texto: %w", err)
	}

	count := 0
	for _, c := range cards {
		if flag == "ativos" && c.Active == 0 {
			continue
		}

		if flag == "inativos" && c.Active == 1 {
			continue
		}

		status := ""
		if flag == "todos" && c.Active == 0 {
			status = "(Inativo)"
		}

		if _, err = sb.WriteString(fmt.Sprintf("#%d - %s %s\n", c.Id, c.Name, status)); err != nil {
			return "", fmt.Errorf("Falha de inclusão de texto: %w", err)
		}
		count++
	}

	if count == 0 {
		return fmt.Sprintf("Nenhum cartão encontrado com o filtro: %s", flag), nil
	}

	return sb.String(), nil
}

func (s *CardService) Update(msg string) (string, error) {
	tokens := strings.Split(msg, " ")
	if len(tokens) < 2 {
		return "", errors.New("Use: [nome] [fnovo_fechamento] [vnovo_vencimento]")
	}

	cMap, err := s.repo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa: %w", err)
	}

	card, exists := cMap[strings.ToLower(tokens[0])]
	if !exists {
		return "", errors.New("Cartão não encontrado. Use /lista c ativos")
	}

	updated := false
	for _, t := range tokens[1:] {
		t = strings.ToLower(t)

		switch {
		case strings.HasPrefix(t, "f"):
			day, err := parseDay(strings.TrimPrefix(t, "f"))
			if err != nil {
				return "", fmt.Errorf("Dia de fechamento inváldio: %w", err)
			}
			card.ClosingDay = day
			updated = true
		case strings.HasPrefix(t, "v"):
			day, err := parseDay(strings.TrimPrefix(t, "v"))
			if err != nil {
				return "", fmt.Errorf("Dia de vencimento inváldio: %w", err)
			}
			card.DueDay = day
			updated = true
		}
	}

	if !updated {
		return "", errors.New("Nenhum parâmetro válido para atualizar. Use f ou v")
	}

	if err = s.repo.Update(&card); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Nome não encontrado. Use /lista c")
		}
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("Cartão (ID %d) atualizada com sucesso! Fechamento: %d e Vencimento: %d", card.Id, card.ClosingDay, card.DueDay), nil
}

func parseDay(s string) (uint8, error) {
	value, err := strconv.ParseUint(s, 10, 8)
	if err != nil || value < 1 || value > 31 {
		return 0, errors.New("use valores de 1 a 31")
	}
	return uint8(value), nil
}

func (s *CardService) Disable(msg string) (string, error) {
	name := strings.Split(strings.TrimSpace(msg), " ")
	if len(name) != 1 {
		return "", errors.New("Use: [nome]")
	}

	cMap, err := s.repo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa: %w", err)
	}

	c, exists := cMap[strings.ToLower(name[0])]
	if !exists {
		return "", errors.New("Cartão não encontrado. Use: /lista c inativos")
	}

	if err = s.repo.Disable(c.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Cartão não encontrado ou já desativado. Use: /lista c inativos")
		}
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("Cartão (ID %d) desativado com sucesso!", c.Id), nil
}

func (s *CardService) Enable(msg string) (string, error) {
	name := strings.Split(strings.TrimSpace(msg), " ")
	if len(name) != 1 {
		return "", errors.New("Use: [nome]")
	}

	cMap, err := s.repo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa: %w", err)
	}

	c, exists := cMap[strings.ToLower(name[0])]
	if !exists {
		return "", errors.New("Cartão não encontrado. Use: /lista c ativos")
	}

	if err = s.repo.Enable(c.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Cartão não encontrado ou já desativado. Use: /lista c ativos")
		}
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("Cartão (ID %d) reativado com sucesso!", c.Id), nil
}

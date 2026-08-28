package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"gfinancer/internal/domain"
	"gfinancer/internal/parser"
	"gfinancer/internal/repository"
)

type ExpenseService struct {
	repo  *repository.ExpenseRepo
	pRepo *repository.PersonRepo
	cRepo *repository.CardRepo
}

func NewExpenseService(r *repository.ExpenseRepo, p *repository.PersonRepo, c *repository.CardRepo) *ExpenseService {
	return &ExpenseService{
		repo:  r,
		pRepo: p,
		cRepo: c,
	}
}

func (s *ExpenseService) Create(msg string) (string, error) {
	pMap, err := s.pRepo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar mapa de pessoas: %w", err)
	}
	cMap, err := s.cRepo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar mapa de cards: %w", err)
	}

	newExpense, card, err := parser.ExpenseParser(msg, pMap, cMap)
	if err != nil {
		return "", fmt.Errorf("Erro de parser: %w", err)
	}

	if err = s.repo.Save(newExpense, card); err != nil {
		return "", fmt.Errorf("Falha ao salvar compra no banco de dados: %w", err)
	}

	return fmt.Sprintf("Compra (ID %d) salva com sucesso!", newExpense.Id), nil
}

func (s *ExpenseService) List(msg string) (string, error) {
	now := time.Now()
	month := int(now.Month())
	year := int(now.Year())
	var personId uint = 0
	var cardId uint = 0

	msg = strings.TrimSpace(msg)
	if msg != "" {
		tokens := strings.Fields(msg)

		if m, err := strconv.Atoi(tokens[0]); err == nil {
			if m < 1 || m > 12 {
				return "", errors.New("Mês inválido. Use de 1 a 12")
			}
			month = m
		} else {
			switch strings.ToLower(tokens[0]) {
			case "p":
				if len(tokens) < 2 {
					return "", errors.New("Use: /lista e p [nome]")
				}
				pMap, err := s.pRepo.ReadMap()
				if err != nil {
					return "", fmt.Errorf("Falha ao carregar o mapa de pessoas: %w", err)
				}
				p, exists := pMap[strings.ToLower(tokens[1])]
				if !exists {
					return "", errors.New("Pessoa não encontrada. Use /lista p")
				}
				personId = p.Id

				if len(tokens) == 3 {
					month, err = strconv.Atoi(tokens[2])
					if err != nil || month < 1 || month > 12 {
						return "", errors.New("MêS inválido. Use de 1 a 12")
					}
				}
			case "c":
				if len(tokens) < 2 {
					return "", errors.New("Use: /lista e c [nome]")
				}
				cMap, err := s.cRepo.ReadMap()
				if err != nil {
					return "", fmt.Errorf("Falha ao carregar o mapa de cartões: %w", err)
				}
				c, exists := cMap[strings.ToLower(tokens[1])]
				if !exists {
					return "", errors.New("Cartão não encontrada. Use /lista c")
				}
				cardId = c.Id

				if len(tokens) == 3 {
					month, err = strconv.Atoi(tokens[2])
					if err != nil || month < 1 || month > 12 {
						return "", errors.New("MêS inválido. Use de 1 a 12")
					}
				}
			default:
				return "", errors.New("Filtro inválido. Use um número de mês, ou 'p [nome]' / 'c [nome]'")
			}
		}
	}

	expenses, err := s.repo.Get(month, year, personId, cardId)
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar do banco de dados: %w", err)
	}

	if len(expenses) == 0 {
		return fmt.Sprintf("Nenhuma compra encontrada para %02d/%02d", month, year), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Compras (%02d/%02d): \n\n", month, year))

	for _, e := range expenses {
		valueStr := fmt.Sprintf("R$ %.2f", e.TotalValue)
		sb.WriteString(fmt.Sprintf("- (%d) %s | %s\n", e.Id, e.Description, valueStr))
	}

	return sb.String(), nil
}

func (s *ExpenseService) Update(msg string) (string, error) {
	tokens := strings.SplitN(strings.TrimSpace(msg), " ", 2)
	if len(tokens) != 2 || strings.TrimSpace(tokens[1]) == "" {
		return "", errors.New("Use: [id] [descricao_nova]")
	}

	id, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil {
		return "", errors.New("Id inválido")
	}

	e := &domain.Expense{
		Id:          uint(id),
		Description: tokens[1],
	}

	if err = s.repo.Update(e); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Id não encontrado. Use /lista e")
		}
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return "Descrição atualizada com suceso!", nil
}

func (s *ExpenseService) Delete(msg string) (string, error) {
	idT := strings.Fields(msg)
	if len(idT) != 1 {
		return "", errors.New("Use: [id]. Para ver ids, use /lista e")
	}

	id, err := strconv.ParseUint(idT[0], 10, 64)
	if err != nil {
		return "", errors.New("Id inválido")
	}
	err = s.repo.Delete(uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Id não encontrado")
		}
		return "", fmt.Errorf("Falha ao remover do banco de dados: %w", err)
	}

	return "Compra removida com sucesso!", nil
}

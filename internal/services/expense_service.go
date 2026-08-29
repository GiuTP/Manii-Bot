package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

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

func (s *ExpenseService) parseFilters(msg string) (int, int, uint, uint, error) {
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
				return 0, 0, 0, 0, errors.New("Mês inválido. Use de 1 a 12")
			}
			month = m
		} else {
			switch strings.ToLower(tokens[0]) {
			case "p":
				if len(tokens) < 2 {
					return 0, 0, 0, 0, errors.New("Use: /lista e p [nome]")
				}
				pMap, err := s.pRepo.ReadMap()
				if err != nil {
					return 0, 0, 0, 0, fmt.Errorf("Falha ao carregar o mapa de pessoas: %w", err)
				}
				p, exists := pMap[strings.ToLower(tokens[1])]
				if !exists {
					return 0, 0, 0, 0, errors.New("Pessoa não encontrada. Use /lista p")
				}
				personId = p.Id

				if len(tokens) == 3 {
					month, err = strconv.Atoi(tokens[2])
					if err != nil || month < 1 || month > 12 {
						return 0, 0, 0, 0, errors.New("MêS inválido. Use de 1 a 12")
					}
				}
			case "c":
				if len(tokens) < 2 {
					return 0, 0, 0, 0, errors.New("Use: /lista e c [nome]")
				}
				cMap, err := s.cRepo.ReadMap()
				if err != nil {
					return 0, 0, 0, 0, fmt.Errorf("Falha ao carregar o mapa de cartões: %w", err)
				}
				c, exists := cMap[strings.ToLower(tokens[1])]
				if !exists {
					return 0, 0, 0, 0, errors.New("Cartão não encontrada. Use /lista c")
				}
				cardId = c.Id

				if len(tokens) == 3 {
					month, err = strconv.Atoi(tokens[2])
					if err != nil || month < 1 || month > 12 {
						return 0, 0, 0, 0, errors.New("MêS inválido. Use de 1 a 12")
					}
				}
			default:
				return 0, 0, 0, 0, errors.New("Filtro inválido. Use um número de mês, ou 'p [nome]' / 'c [nome]'")
			}
		}
	}

	return month, year, personId, cardId, nil
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func (s *ExpenseService) List(msg string) (string, error) {
	month, year, personId, cardId, err := s.parseFilters(msg)
	if err != nil {
		return "", fmt.Errorf("%w", err)
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
		valueStr := fmt.Sprintf("R$ %.2f", e.InstallmentValue)

		installStr := ""
		if e.TotalInstallments > 1 {
			installStr = fmt.Sprintf(" (%d/%d)", e.CurrentInstallment, e.TotalInstallments)
		}
		sb.WriteString(fmt.Sprintf("- %s%s | %s\n", e.Description, installStr, valueStr))
	}

	return sb.String(), nil
}

func (s *ExpenseService) GetReportData(msg string) (string, []domain.Report, string, error) {
	month, year, personId, cardId, err := s.parseFilters(msg)
	if err != nil {
		return "", nil, "", err
	}

	expenses, err := s.repo.Get(month, year, personId, cardId)
	if err != nil {
		return "", nil, "", fmt.Errorf("Falha ao busca no banco de dados: %w", err)
	}
	if len(expenses) == 0 {
		return "", nil, "", fmt.Errorf("Nenhuma fatura encontrada para %02d/%04d", month, year)
	}

	pMap, _ := s.pRepo.ReadMap()
	cMap, _ := s.cRepo.ReadMap()

	pById := make(map[uint]string)
	for name, p := range pMap {
		pById[p.Id] = name
	}

	cById := make(map[uint]string)
	for name, c := range cMap {
		cById[c.Id] = name
	}

	var reports []domain.Report
	var total float64

	for _, e := range expenses {
		total += float64(e.InstallmentValue)

		installStr := ""
		if e.TotalInstallments > 1 {
			installStr = fmt.Sprintf(" (%d/%d)", e.CurrentInstallment, e.TotalInstallments)
		}

		personName := "N/A"
		if e.PersonId != nil {
			personName = capitalize(pById[*e.PersonId])
		}

		cardName := "N/A"
		if e.CardId != nil {
			cardName = capitalize(cById[*e.CardId])
		}

		reports = append(reports, domain.Report{
			Date:        e.DueDate,
			Description: e.Description,
			Value:       fmt.Sprintf("R$ %.2f", e.InstallmentValue),
			Installment: installStr,
			Person:      personName,
			Card:        cardName,
		})
	}

	title := fmt.Sprintf("EXTRATO - %02d/%04d", month, year)
	totalStr := fmt.Sprintf("R$ %.2f", total)

	return title, reports, totalStr, nil
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

	return "Descrição atualizada com sucesso!", nil
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

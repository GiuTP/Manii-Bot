package services

import (
	"errors"
	"gfinancer/internal/domain"
	"gfinancer/internal/parser"
	"gfinancer/internal/repository"
	"strings"
)

type BotService struct {
	expenseRepo *repository.ExpenseRepo
	cardRepo    *repository.CardRepo
	personRepo  *repository.PersonRepo
}

func NewBotService(e *repository.ExpenseRepo, c *repository.CardRepo, p *repository.PersonRepo) *BotService {
	return &BotService{
		expenseRepo: e,
		cardRepo:    c,
		personRepo:  p,
	}
}

func (s *BotService) HandleMessage(rawMsg string) (string, error) {
	tokens := strings.SplitN(rawMsg, " ", 2)
	cmd := tokens[0]

	var msg string
	if len(tokens) > 1 {
		msg = tokens[1]
	}

	switch cmd {
	case "/expense", "e":
	case "/person", "p":
	case "/card", "c":
	default:
	}
}

func (s *BotService) createExpense(message string) (*domain.Expense, error) {
	personsMap, err := s.personRepo.ReadMap()
	if err != nil {
		return nil, errors.New("erro ao carregar pessoas do banco: " + err.Error())
	}

	cardsMap, err := s.cardRepo.ReadMap()
	if err != nil {
		return nil, errors.New("erro ao carregar cartões do banco: " + err.Error())
	}

	expense, err := parser.ExpenseParser(message, personsMap, cardsMap)
	if err != nil {
		return nil, errors.New("erro ao processar texto: " + err.Error())
	}

	var selectedCard *domain.Card
	if expense.CardId != nil {
		for _, c := range cardsMap {
			if c.Id == *expense.CardId {
				cardCopy := c
				selectedCard = &cardCopy
				break
			}
		}
	}

	err = s.expenseRepo.Save(expense, selectedCard)
	if err != nil {
		return nil, errors.New("erro ao salvar no banco de dados: " + err.Error())
	}

	return expense, nil
}

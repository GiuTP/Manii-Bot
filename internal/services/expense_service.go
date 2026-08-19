package services

import (
	"errors"
	"gfinancer/internal/domain"
	"gfinancer/internal/parser"
	"gfinancer/internal/repository"
)

type ExpenseService struct {
	expenseRepo *repository.ExpenseRepo
	cardRepo    *repository.CardRepo
	personRepo  *repository.PersonRepo
}

func NewExpenseService(e *repository.ExpenseRepo, c *repository.CardRepo, p *repository.PersonRepo) *ExpenseService {
	return &ExpenseService{
		expenseRepo: e,
		cardRepo:    c,
		personRepo:  p,
	}
}

func (s *ExpenseService) ProcessTransaction(message string) (*domain.Expense, error) {
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

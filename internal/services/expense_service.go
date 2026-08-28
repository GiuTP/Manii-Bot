package services

import "gfinancer/internal/repository"

type ExpenseService struct {
	repo *repository.ExpenseRepo
}

func NewExpenseService(r *repository.ExpenseRepo) *ExpenseService {
	return &ExpenseService{repo: r}
}

func (s *ExpenseService) Create(msg string) (string, error) {
	
}

func (s *ExpenseService) List(msg string) (string, error) {

}

func (s *ExpenseService) Update(msg string) (string, error) {

}

func (s *ExpenseService) Delete(msg string) (string, error) {

}

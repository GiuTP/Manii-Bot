package parser

import (
	"gfinancer/internal/domain"
	"testing"
)

func TestExpenseParser(t *testing.T) {
	mockPeople := map[string]domain.Person{
		"giuliano": {Id: 1},
		"giuliany": {Id: 2},
		"patricia": {Id: 3},
		"melissa":  {Id: 4},
		"aylton":   {Id: 5},
		"ryan":     {Id: 6},
		"carina":   {Id: 7},
		"henrique": {Id: 8},
	}

	mockCards := map[string]domain.Card{
		"bradesco": {Id: 1},
		"mp":       {Id: 2},
	}

	tests := []struct {
		name             string
		input            string
		waitError        bool
		waitValue        float64
		waitDesc         string
		waitInstallments uint8
		waitPerson       uint
		waitCard         uint
	}{
		{
			name:             "Gasto básico apenas com valor e descrição",
			input:            "45.90 Ifood",
			waitError:        false,
			waitValue:        45.90,
			waitDesc:         "Ifood",
			waitInstallments: 1,
		},
		{
			name:             "Gasto completo com vírgula, fora de ordem",
			input:            "150,50 burger king giuliany bradesco 3x",
			waitError:        false,
			waitValue:        150.50,
			waitDesc:         "burger king",
			waitInstallments: 3,
			waitPerson:       2,
			waitCard:         1,
		},
		{
			name:             "Gasto com outra ordem e outro cartao",
			input:            "20.00 uber 2x mp giuliano",
			waitError:        false,
			waitValue:        20.00,
			waitDesc:         "uber",
			waitInstallments: 2,
			waitPerson:       1,
			waitCard:         2,
		},
		{
			name:      "Erro: Falta a descrição",
			input:     "50.00 mp",
			waitError: true,
		},
		{
			name:      "Erro: Valor inválido",
			input:     "abc Ifood",
			waitError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expense, _, err := ExpenseParser(tt.input, mockPeople, mockCards)

			if tt.waitError {
				if err == nil {
					t.Errorf("Esperava um erro, mas não ocorreu")
				}
				return
			}

			if err != nil {
				t.Fatalf("Não esperava erro, mas ocorreu: %v", err)
			}

			if expense.TotalValue != tt.waitValue {
				t.Errorf("Valor incorreto. Esperado: %v, Recebido: %v", tt.waitValue, expense.TotalValue)
			}
			if expense.Description != tt.waitDesc {
				t.Errorf("Descrição incorreta. Esperada: '%v', Recebida: '%v'", tt.waitDesc, expense.Description)
			}
			if expense.TotalInstallments != tt.waitInstallments {
				t.Errorf("Parcelas incorretas. Esperado: %v, Recebido: %v", tt.waitInstallments, expense.TotalInstallments)
			}

			if tt.waitPerson != 0 {
				if expense.PersonId == nil || *expense.PersonId != tt.waitPerson {
					t.Errorf("Pessoa incorreta. Esperada: %v", tt.waitPerson)
				}
			}

			if tt.waitCard != 0 {
				if expense.CardId == nil || *expense.CardId != tt.waitCard {
					t.Errorf("Cartão incorreto. Esperado: %v", tt.waitCard)
				}
			}
		})
	}
}

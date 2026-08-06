package domain

import (
	"math"
	"time"
)

// InstallmentGenarate cria novas compras menores de um cartão (card) a partir de uma compra caso esta tenha sido parcelada.
// Retorna um slice de Installment com os valores dos parcelamentos e error.
func (g *Expense) InstallmentGenerate(card *Card) ([]Installment, error) {
	// Slice de parcelas que será retornado
	var installments []Installment

	// Verificação da formatação da data de compra
	PurchaseDate, err := time.Parse("2006-01-02", g.PurchaseDate)
	if err != nil {
		return nil, err
	}

	baseDate := PurchaseDate

	// Cartão existe e é cartão de crédito
	if card != nil && card.Type == 0 {
		// Valores são lançados conforme regras de cartões de crédito
		if PurchaseDate.Day() >= int(card.ClosingDay) {
			baseDate = baseDate.AddDate(0, 1, 0)
		}
		if card.DueDay < card.ClosingDay {
			baseDate = baseDate.AddDate(0, 1, 0)
		}

		// Data formatada como YYYY-MM-DD
		baseDate = time.Date(baseDate.Year(), baseDate.Month(), int(card.DueDay), 0, 0, 0, 0, baseDate.Location())
	}

	// Valor base do parcelamento
	// Pega apenas duas casas decimais de números reais
	baseValue := math.Floor((g.TotalValue/float64(g.TotalInstallments))*100) / 100

	// Problemas de centavos
	// Compras com dizimas, a primeira parcela absorve a diferença
	firstInstallmentValue := g.TotalValue - (baseValue * float64(g.TotalInstallments-1))
	firstInstallmentValue = math.Round(firstInstallmentValue*100) / 100

	// Inclui os valores no slice de parcelas
	// Os valores incluídos serão o valor base ou o valor da primeira parcela.
	firstInstallment := Installment{
		ExpenseId:          g.Id,
		NumberInstallments: 1,
		Value:              firstInstallmentValue,
		DueDate:            baseDate.Format("2006-01-02"),
		PaymentStatus:      0,
	}
	installments = append(installments, firstInstallment)

	for i := uint8(2); i <= g.TotalInstallments; i++ {
		installment := Installment{
			ExpenseId:          g.Id,
			NumberInstallments: i,
			Value:              baseValue,
			DueDate:            baseDate.AddDate(0, int(i-1), 0).Format("2006-01-02"),
			PaymentStatus:      0,
		}
		installments = append(installments, installment)
	}

	return installments, nil
}

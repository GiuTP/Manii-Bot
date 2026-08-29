package domain

// Expense representa o detalhamento de uma compra.
type Expense struct {
	Id                uint    // Identificação
	Description       string  // Descrição da compra
	TotalValue        float64 // Valor total da compra
	PurchaseDate      string  // Dia da compra
	TotalInstallments uint8   // Quantidade de parcelas
	PersonId          *uint   // Quem fez a compra. Campo opcional
	CardId            *uint   // Cartão onde foi feito a compra. Campo opcional

	InstallmentValue   float64 // Valor da parcela atual
	CurrentInstallment uint8   // Parcela atual
	DueDate            string  // Data de vencimento da parcela
}

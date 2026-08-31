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
}

// ExpenseDTO é usada para transferir os dados mesclados de despesas e parcelas para a listagem
type ExpenseDTO struct {
	Expense
	InstallmentValue   float64
	CurrentInstallment uint8
	DueDate            string
}

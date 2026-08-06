package domain

// Person representa quem fez a compra. Identificada em Expense.
type Person struct {
	Id   uint   // Identificação
	Name string // Nome da pessoa
}

// Card representa onde foi feito a compra. Identificado em Expense.
type Card struct {
	Id         uint   // Identificação
	Name       string // Nome do cartão (banco)
	Type       uint8  // Tipo do cartão. 0 para crédito e 1 para débito
	ClosingDay uint8  // Dia de fechamento da fatura do cartão
	DueDay     uint8  // Dia de vencimento da fatura do cartão
}

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

// Installment representa o parcelamento de uma compra. Usa Expense.
type Installment struct {
	Id                 uint    // Idenficação
	ExpenseId          uint    // Compra original em Expense
	NumberInstallments uint8   // Parcela atual. Usado para mostrar em qual parcela estar.
	Value              float64 // Valor da parcela
	DueDate            string  // Data de vencimento
	PaymentStatus      uint8   // Status de parcela. 0 para pendente e 1 para paga.
}

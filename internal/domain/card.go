package domain

// Card representa onde foi feito a compra. Identificado em Expense.
type Card struct {
	Id         uint   // Identificação
	Name       string // Nome do cartão (banco)
	Active     uint8  // Flag de existência de card. 0 inativo, 1 ativo
	Type       uint8  // Tipo do cartão. 0 para crédito e 1 para débito
	ClosingDay uint8  // Dia de fechamento da fatura do cartão
	DueDay     uint8  // Dia de vencimento da fatura do cartão
}

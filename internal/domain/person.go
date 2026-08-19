package domain

// Person representa quem fez a compra. Identificada em Expense.
type Person struct {
	Id     uint   // Identificação
	Name   string // Nome da pessoa
	Active uint8  // Flag de existência de person. 0 para inativo e 1 para ativo
}

package domain

type Pessoa struct {
	Id   uint
	Nome string
}

type Cartao struct {
	Id            uint
	Nome          string
	Tipo          uint8
	DiaFechamento uint8
	DiaVencimento uint8
}

type Gasto struct {
	Id            uint
	Descricao     string
	ValorTotal    float64
	DataCompra    string
	TotalParcelas uint8
	PessoaId      *uint
	CartaoId      *uint
}

type Parcela struct {
	Id              uint
	GastoId         uint
	NumeroParcela   uint8
	Valor           float64
	DataCobranca    string
	StatusPagamento uint8
}

package domain

import (
	"math"
	"time"
)

func (g *Gasto) GerarParcelas(cartao *Cartao) ([]Parcela, error) {
	var parcelas []Parcela

	dataCompra, err := time.Parse("2006-01-02", g.DataCompra)
	if err != nil {
		return nil, err
	}

	dataBase := dataCompra

	if cartao != nil && cartao.Tipo == 0 {
		if dataCompra.Day() >= int(cartao.DiaFechamento) {
			dataBase = dataBase.AddDate(0, 1, 0)
		}

		if cartao.DiaVencimento < cartao.DiaFechamento {
			dataBase = dataBase.AddDate(0, 1, 0)
		}

		dataBase = time.Date(dataBase.Year(), dataBase.Month(), int(cartao.DiaVencimento), 0, 0, 0, 0, dataBase.Location())
	}

	valorBase := math.Floor((g.ValorTotal / float64(g.TotalParcelas)) * 100 / 100)

	primeiraParcelaValor := g.ValorTotal - (valorBase * float64(g.TotalParcelas-1))
	primeiraParcelaValor = math.Round(primeiraParcelaValor*100) / 100

	for i := uint8(1); i < g.TotalParcelas; i++ {
		dataCobranca := dataBase.AddDate(0, int(i-1), 0)
		valorParcela := valorBase
		if i == 1 {
			valorParcela = primeiraParcelaValor
		}

		parcela := Parcela{
			GastoId:         g.Id,
			NumeroParcela:   i,
			Valor:           valorParcela,
			DataCobranca:    dataCobranca.Format("2006-01-02"),
			StatusPagamento: 0,
		}

		parcelas = append(parcelas, parcela)
	}

	return parcelas, nil
}

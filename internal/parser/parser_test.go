package parser

import (
	"testing"
)

func TestParseGasto(t *testing.T) {
	mockPessoas := map[string]uint{
		"giuliano": 1,
		"giuliany": 2,
		"patricia": 3,
		"melissa":  4,
		"aylton":   5,
		"ryan":     6,
		"carina":   7,
		"henrique": 8,
	}

	mockCartoes := map[string]uint{
		"bradesco": 1,
		"mp":       2,
	}

	testes := []struct {
		nome           string
		input          string
		esperaErro     bool
		esperaValor    float64
		esperaDesc     string
		esperaParcelas uint8
		esperaPessoa   uint
		esperaCartao   uint
	}{
		{
			nome:           "Gasto básico apenas com valor e descrição",
			input:          "45.90 Ifood",
			esperaErro:     false,
			esperaValor:    45.90,
			esperaDesc:     "Ifood",
			esperaParcelas: 1, // Padrão
		},
		{
			nome:           "Gasto completo com vírgula, fora de ordem",
			input:          "150,50 burger king giuliany bradesco 3x",
			esperaErro:     false,
			esperaValor:    150.50,
			esperaDesc:     "burger king", // Juntou o que sobrou
			esperaParcelas: 3,
			esperaPessoa:   2, // ID da giuliany no mock
			esperaCartao:   1, // ID do bradesco no mock
		},
		{
			nome:           "Gasto com outra ordem e outro cartao",
			input:          "20.00 uber 2x bradesco giuliano",
			esperaErro:     false,
			esperaValor:    20.00,
			esperaDesc:     "uber",
			esperaParcelas: 2,
			esperaPessoa:   1,
			esperaCartao:   1,
		},
		{
			nome:       "Erro: Falta a descrição",
			input:      "50.00 mp",
			esperaErro: true, // "mp" é removido e não sobra nada para a descrição
		},
		{
			nome:       "Erro: Valor inválido",
			input:      "abc Ifood",
			esperaErro: true,
		},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			gasto, err := ParseGasto(tt.input, mockPessoas, mockCartoes)

			if tt.esperaErro {
				if err == nil {
					t.Errorf("Esperava um erro, mas não ocorreu")
				}
				return // Se esperava erro e deu erro, o teste passou, vai pro próximo
			}

			if err != nil {
				t.Fatalf("Não esperava erro, mas ocorreu: %v", err)
			}

			if gasto.ValorTotal != tt.esperaValor {
				t.Errorf("Valor incorreto. Esperado: %v, Recebido: %v", tt.esperaValor, gasto.ValorTotal)
			}
			if gasto.Descricao != tt.esperaDesc {
				t.Errorf("Descrição incorreta. Esperada: '%v', Recebida: '%v'", tt.esperaDesc, gasto.Descricao)
			}
			if gasto.TotalParcelas != tt.esperaParcelas {
				t.Errorf("Parcelas incorretas. Esperado: %v, Recebido: %v", tt.esperaParcelas, gasto.TotalParcelas)
			}

			// Verificações com ponteiros exigem cuidado para não dar erro de nil pointer
			if tt.esperaPessoa != 0 {
				if gasto.PessoaId == nil || *gasto.PessoaId != tt.esperaPessoa {
					t.Errorf("Pessoa incorreta. Esperada: %v", tt.esperaPessoa)
				}
			}

			if tt.esperaCartao != 0 {
				if gasto.CartaoId == nil || *gasto.CartaoId != tt.esperaCartao {
					t.Errorf("Cartão incorreto. Esperado: %v", tt.esperaCartao)
				}
			}
		})
	}
}

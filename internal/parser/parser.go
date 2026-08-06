package parser

import (
	"errors"
	"gfinancer/internal/domain"
	"regexp"
	"strconv"
	"strings"
)

func ParseGasto(input string, pessoas map[string]uint, cartoes map[string]uint) (*domain.Gasto, error) {
	input = strings.TrimSpace(input)

	// Input vazio é inválido
	if input == "" {
		return nil, errors.New("input vazio")
	}

	tokens := strings.Fields(input)

	// Input abaixo de 2 tokens é inválido
	// É esperado <valor> <descrição> pelo menos
	if len(tokens) < 2 {
		return nil, errors.New("formato invalido, esperado pelo menos: <valor> <descricao>")
	}

	gasto := &domain.Gasto{
		TotalParcelas: 1,
	}

	valorStr := strings.Replace(tokens[0], ",", ".", 1)
	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil {
		return nil, errors.New("o primeiro termo deve ser um valor numerico")
	}
	gasto.ValorTotal = valor

	reParcela := regexp.MustCompile(`^(\d+)[xX]$`)

	var descTokens []string

	for _, token := range tokens[1:] {
		tokenLower := strings.ToLower(token)

		if match := reParcela.FindStringSubmatch(tokenLower); match != nil {
			parcelas, _ := strconv.Atoi(match[1])
			gasto.TotalParcelas = uint8(parcelas)
			continue
		}

		if id, existe := pessoas[tokenLower]; existe {
			idCopy := id
			gasto.PessoaId = &idCopy
			continue
		}

		if id, existe := cartoes[tokenLower]; existe {
			idCopy := id
			gasto.CartaoId = &idCopy
			continue
		}

		descTokens = append(descTokens, token)
	}

	gasto.Descricao = strings.Join(descTokens, " ")

	if gasto.Descricao == "" {
		return nil, errors.New("a descricao nao pode ser vazia")
	}

	return gasto, nil
}

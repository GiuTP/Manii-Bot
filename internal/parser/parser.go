package parser

import (
	"errors"
	"gfinancer/internal/domain"
	"regexp"
	"strconv"
	"strings"
)

func ExpenseParser(input string, people map[string]uint, cards map[string]uint) (*domain.Expense, error) {
	input = strings.TrimSpace(input)

	// Input vazio é inválido
	if input == "" {
		return nil, errors.New("input vazio")
	}

	// Cria os tokens a partir de input
	tokens := strings.Fields(input)

	// Input abaixo de 2 tokens é inválido
	// É esperado <valor> <descrição> pelo menos
	if len(tokens) < 2 {
		return nil, errors.New("formato invalido, esperado pelo menos: <valor> <descricao>")
	}

	// Numéro de parcelas é 1 por padrão
	// Usado caso não seja informado o número.
	expense := &domain.Expense{
		TotalInstallments: 1,
	}

	// Converte o padrão brasileiro de decimal para o americano
	valueStr := strings.Replace(tokens[0], ",", ".", 1)
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, errors.New("o primeiro termo deve ser um valor numerico")
	}
	expense.TotalValue = value

	// Regex para procurar token de parcelamento
	reInstallment := regexp.MustCompile(`^(\d+)[xX]$`)

	var descTokens []string

	// Varre os tokens restantes
	// Procura por: parcelamento, pessoa, cartão e descrição
	for _, token := range tokens[1:] {
		tokenLower := strings.ToLower(token)

		// Parcelamento
		if match := reInstallment.FindStringSubmatch(tokenLower); match != nil {
			installments, _ := strconv.Atoi(match[1])
			expense.TotalInstallments = uint8(installments)
			continue
		}

		// Pessoa
		if id, existe := people[tokenLower]; existe {
			idCopy := id
			expense.PersonId = &idCopy
			continue
		}

		// Cartão
		if id, existe := cards[tokenLower]; existe {
			idCopy := id
			expense.CardId = &idCopy
			continue
		}

		// Descrição (resto)
		descTokens = append(descTokens, token)
	}

	expense.Description = strings.Join(descTokens, " ")
	if expense.Description == "" {
		return nil, errors.New("a descricao nao pode ser vazia")
	}

	return expense, nil
}

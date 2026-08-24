package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"gfinancer/internal/domain"
)

// ExpenseParser converte uma mensagem de texto puro em uma estrutura de despesa.
// A função extrai o valor monetário, número de parcelas, pessoa e cartão
// com base nos mapas fornecidos, unindo o restante das palavras como descrição.
// Retorna a despesa preenchida em caso de sucesso, ou um erro se o texto
// for inválido ou não contiver um valor numérico.
func ExpenseParser(input string, people map[string]domain.Person, cards map[string]domain.Card) (*domain.Expense, error) {
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
		return nil, errors.New("formato inválido, esperado pelo menos descrição e valor da compra.")
	}

	// Numéro de parcelas é 1 por padrão
	// Usado caso não seja informado o número.
	expense := &domain.Expense{
		TotalInstallments: 1,
	}

	// Regex para procurar token de parcelamento
	reInstallment := regexp.MustCompile(`^(\d+)[xX]$`)
	// Regex para procurar token de valor (dinheiro)
	reMoney := regexp.MustCompile(`^\d+[.,]\d{2}$`)

	var descTokens []string
	var valueFound bool

	// Varre os tokens restantes
	// Procura por: parcelamento, pessoa, cartão e descrição
	for _, token := range tokens {
		tokenLower := strings.ToLower(token)

		// Valor
		if !valueFound && reMoney.MatchString(tokenLower) {
			valueStr := strings.Replace(tokenLower, ",", ".", 1)
			if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
				expense.TotalValue = value
				valueFound = true
				continue
			}
		}

		// Parcelamento
		if match := reInstallment.FindStringSubmatch(tokenLower); match != nil {
			installments, _ := strconv.Atoi(match[1])
			expense.TotalInstallments = uint8(installments)
			continue
		}

		// Pessoa
		if p, existe := people[tokenLower]; existe {
			pCopy := p
			expense.PersonId = &pCopy.Id
			continue
		}

		// Cartão
		if c, existe := cards[tokenLower]; existe {
			cCopy := c
			expense.CardId = &cCopy.Id
			continue
		}

		// Descrição (resto)
		descTokens = append(descTokens, token)
	}

	// Valor não encontrado
	// Caso de uso: compra grátis.
	if !valueFound {
		return nil, errors.New("valor não encontrado. Valor deve conter casas decimais, mesmo que '.00'")
	}

	// Une tokens de descrição para casos onde contém mais de uma string (mais comum)
	expense.Description = strings.Join(descTokens, " ")
	if expense.Description == "" {
		return nil, errors.New("a descrição não pode ser vazia")
	}

	return expense, nil
}

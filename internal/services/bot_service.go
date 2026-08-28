package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gfinancer/internal/domain"
	"gfinancer/internal/parser"
)

type BotService struct {
	eSvc *ExpenseService
	cSvc *CardService
	pSvc *PersonService
}

func NewBotService(e *ExpenseService, c *CardService, p *PersonService) *BotService {
	return &BotService{
		eSvc: e,
		cSvc: c,
		pSvc: p,
	}
}

func (s *BotService) HandleMessage(rawMsg string) (string, error) {
	tokens := strings.SplitN(rawMsg, " ", 2)
	if len(tokens) < 1 {
		return "", errors.New("Comando vazio. Tente /help")
	}
	cmd := tokens[0]
	msg := tokens[1]

	// Todos os comandos do bot
	switch cmd {
	case "/compra":
		exp, err := s.createExpense(msg)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Despesa de R$ %.2f registrada.", exp.TotalValue), nil
	case "/pessoa":
		p, err := s.createPerson(msg)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s registrado(a).", p.Name), nil
	case "/cartao":
		c, err := s.createCard(msg)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s registrado(a).", c.Name), nil
	case "/listar":
		l, err := s.list(msg)
		if err != nil {
			return "", err
		}
		return l, nil
	case "/apagar":
	case "/atualizar":
		err := s.update(msg)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Atualizado com sucesso"), nil
	case "/help":
		return s.listCmds(), nil
	default:
		return "Comando inexistente. Para listagem digite /help", nil
	}

	// Apenas para o compilador não reclamar
	return "", nil
}

// ----------------------------------------------------------------------------
// - Funções de atualização
// ----------------------------------------------------------------------------

func (s *BotService) update(msg string) error {
	tokens := strings.SplitN(strings.TrimSpace(msg), " ", 2)
	if len(tokens) < 2 {
		return errors.New("formato inválido.")
	}

	switch tokens[0] {
	case "p":
		err := s.updatePerson(tokens[1])
		if err != nil {
			return err
		}
	default:
		return errors.New("Sub comando inválido. Tente p, c ou e;")
	}

	return nil
}

// ----------------------------------------------------------------------------
// - Funções de exclusão/desativação
// ----------------------------------------------------------------------------
func (s *BotService) delete(msg string) error {
	tokens := strings.Split(strings.TrimSpace(msg), " ")

	if len(tokens) != 2 {
		return errors.New("Use: [subcomando] [id]")
	}

	switch tokens[0] {
	case "p":
		return s.disablePerson(tokens[1])
	case "c":
		return s.disableCard(tokens[1])
	case "e":
		return s.deleteExpense(tokens[1])
	default:
		return errors.New("subcomando não existe. Tente p, c ou e")
	}
}

func (s *BotService) disableCard(c string) error {
	cMap, err := s.cardRepo.ReadMap()
	if err != nil {
		return err
	}
	dCard, exist := cMap[c]
	if !exist {
		return errors.New("cartão não existe")
	}
	if err = s.cardRepo.Disable(dCard.Id); err != nil {
		return err
	}

	return nil
}

func (s *BotService) deleteExpense(e string) error {
	eId, err := strconv.Atoi(e)
	if err != nil {
		return errors.New("valor precisa ser um número")
	}
	if err = s.expenseRepo.Delete(uint(eId)); err != nil {
		return err
	}
	return nil
}

func (s *BotService) list(msg string) (string, error) {
	token := strings.TrimSpace(msg)

	switch token {
	case "p":
		p, err := s.listPerson()
		if err != nil {
			return "", err
		}
		return p, nil
	// case "c":
	// case "e":
	default:
		return "", errors.New("subcomando inexistente. Tente p, c ou e")
	}

	// return "", nil
}

// Melhorar
func (s *BotService) listCmds() string {
	return `*Comandos disponíveis:*
	- Nova despesa
	/expense | /e | /compra [msg]

	- Nova cessoa
	/person | /p | /pessoa [nome]

	- Novo cartão
	/card | /c | /cartao | /cartão [nome] [tipo] [fech] [venc]

	- Deletar (BREVE)
	/delete | /d [id]
	`
}

func (s *BotService) createCard(msg string) (*domain.Card, error) {
	// Validação de formato
	tokens := strings.Split(strings.TrimSpace(msg), " ")
	if len(tokens) != 4 {
		return nil, errors.New("Use: [Nome] [Tipo] [Dia Fechamento] [Dia Vencimento]")
	}

	name := tokens[0]

	// Validação de tipo do cartão
	var typ uint8
	switch tokens[1] {
	case "crédito", "credito":
		typ = 0
	case "débito", "debito":
		typ = 1
	default:
		return nil, errors.New("tipo do cartão inválido")
	}

	closingD, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, errors.New("dia de fechamento deve ser um número")
	}
	dueD, err := strconv.Atoi(tokens[3])
	if err != nil {
		return nil, errors.New("dia de vencimento deve ser um número")
	}

	newC := &domain.Card{
		Name:       name,
		Type:       typ,
		ClosingDay: uint8(closingD),
		DueDay:     uint8(dueD),
	}
	err = s.cardRepo.Create(newC)
	if err != nil {
		return nil, err
	}

	return newC, nil
}

func (s *BotService) createExpense(message string) (*domain.Expense, error) {
	personsMap, err := s.personRepo.ReadMap()
	if err != nil {
		return nil, errors.New("erro ao carregar pessoas do banco: " + err.Error())
	}

	cardsMap, err := s.cardRepo.ReadMap()
	if err != nil {
		return nil, errors.New("erro ao carregar cartões do banco: " + err.Error())
	}

	expense, err := parser.ExpenseParser(message, personsMap, cardsMap)
	if err != nil {
		return nil, errors.New("erro ao processar texto: " + err.Error())
	}

	var selectedCard *domain.Card
	if expense.CardId != nil {
		for _, c := range cardsMap {
			if c.Id == *expense.CardId {
				cardCopy := c
				selectedCard = &cardCopy
				break
			}
		}
	}

	err = s.expenseRepo.Save(expense, selectedCard)
	if err != nil {
		return nil, errors.New("erro ao salvar no banco de dados: " + err.Error())
	}

	return expense, nil
}

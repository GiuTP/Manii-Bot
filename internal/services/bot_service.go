package services

import (
	"errors"
	"strings"
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

func (s *BotService) HandleMessage(rawMsg string) string {
	tokens := strings.SplitN(rawMsg, " ", 2)
	if len(tokens) < 1 {
		return "Comando vazio."
	}

	cmd := tokens[0]
	msg := ""
	if len(tokens) == 2 {
		msg = strings.TrimSpace(tokens[1])
	}

	var resp string
	var err error

	// Todos os comandos do bot
	switch cmd {

	// Criação
	case "/compra":
		resp, err = s.eSvc.Create(msg)
	case "/pessoa":
		resp, err = s.pSvc.Create(msg)
	case "/cartao":
		resp, err = s.cSvc.Create(msg)

	// Listagem
	case "/listar":
		resp, err = s.list(msg)

	//Atualização
	case "/atualizar":
		resp, err = s.update(msg)

	// Exclusão
	case "/apagar":
		resp, err = s.delete(msg)

	// Reativação
	case "/ativar":

	// Help
	case "/help":

	default:
		return "Comando inexistente. Para listagem digite /help"
	}

	if err != nil {
		return "(Falha) " + err.Error()
	}

	return "(Sucesso) " + resp
}

func (s *BotService) update(msg string) (string, error) {
	if msg = strings.TrimSpace(msg); msg == "" {
		return "", errors.New("Especifique qual entidade atualizar. Use: p, c ou e")
	}

	tokens := strings.SplitN(msg, " ", 2)
	if len(tokens) < 2 {
		return "", errors.New("Use: [entidade] [id/nome] [novos_valores]")
	}

	cmd := strings.ToLower(tokens[0])
	upd := strings.TrimSpace(tokens[1])

	switch cmd {
	case "p":
		return s.pSvc.Update(upd)
	case "c":
		return s.cSvc.Update(upd)
	case "e":
		return s.eSvc.Update(upd)
	default:
		return "", errors.New("Subcomando inexistente. Use: p, c ou e")
	}
}

func (s *BotService) delete(msg string) (string, error) {
	if msg = strings.TrimSpace(msg); msg == "" {
		return "", errors.New("Especifique qual entidade apagar/desativar. Use: p, c ou e")
	}

	tokens := strings.SplitN(msg, " ", 2)
	if len(tokens) < 2 {
		return "", errors.New("Use: [entidade] [id/nome]")
	}
	cmd := strings.ToLower(tokens[0])
	del := strings.TrimSpace(tokens[1])

	switch cmd {
	case "p":
		return s.pSvc.Disable(del)
	case "c":
		return s.cSvc.Disable(del)
	case "e":
		return s.eSvc.Delete(del)
	default:
		return "", errors.New("Subcomando inexistente. Use: p, c ou e")
	}
}

func (s *BotService) list(msg string) (string, error) {
	if msg = strings.TrimSpace(msg); msg == "" {
		return "", errors.New("Especifique qual entidade listar. Use: p, c ou e")
	}

	tokens := strings.SplitN(msg, " ", 2)
	cmd := strings.ToLower(tokens[0])
	flag := ""
	if len(tokens) == 2 {
		flag = strings.ToLower(strings.TrimSpace(tokens[1]))
	}

	switch cmd {
	case "p":
		return s.pSvc.List(flag)
	case "c":
		return s.cSvc.List(flag)
	case "e":
		return s.eSvc.List(flag)
	default:
		return "", errors.New("Subcomando inexistente. Use: p, c ou e")
	}
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

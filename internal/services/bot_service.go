package services

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"
)

type BotService struct {
	eSvc *ExpenseService
	cSvc *CardService
	pSvc *PersonService
	rSvc *ReportService
	sSvc *SubscriptionService
}

func NewBotService(e *ExpenseService, c *CardService, p *PersonService, r *ReportService, s *SubscriptionService) *BotService {
	return &BotService{
		eSvc: e,
		cSvc: c,
		pSvc: p,
		rSvc: r,
		sSvc: s,
	}
}

func (s *BotService) HandleMessage(rawMsg string) (string, string) {
	tokens := strings.SplitN(rawMsg, " ", 2)
	if len(tokens) < 1 {
		return "Comando vazio.", ""
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
		resp, err = s.enable(msg)

	// Relatório
	case "/relatorio":
		title, reports, totalStr, err := s.eSvc.GetReportData(msg)
		if err != nil {
			return fmt.Sprintf("Erro: %v", err), ""
		}

		filePath, err := s.rSvc.GeneratePDF(title, reports, totalStr)
		if err != nil {
			return fmt.Sprintf("Erro ao gerar arquivo PDF: %v", err), ""
		}

		return "Aqui está seu extrato!", filePath

	// Help
	case "/help":
		resp = s.help()

	// Comando incorreto
	default:
		return "Comando inexistente. Para listagem digite /help", ""
	}

	if err != nil {
		return "(Falha) " + err.Error(), ""
	}

	return "(Sucesso) " + resp, ""
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
	case "s":
		return s.sSvc.Update(upd)
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
	case "s":
		return s.sSvc.Cancel(del)
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
	case "s":
		return s.sSvc.List(flag)
	default:
		return "", errors.New("Subcomando inexistente. Use: p, c ou e")
	}
}

func (s *BotService) enable(msg string) (string, error) {
	if msg = strings.TrimSpace(msg); msg == "" {
		return "", errors.New("Especifique qual entidade reativar. Use: p ou c")
	}

	tokens := strings.SplitN(msg, " ", 2)
	if len(tokens) < 2 {
		return "", errors.New("Use: [entidade] [nome]")
	}
	cmd := strings.ToLower(tokens[0])
	act := strings.TrimSpace(tokens[1])

	switch cmd {
	case "p":
		return s.pSvc.Enable(act)
	case "c":
		return s.cSvc.Enable(act)
	default:
		return "", errors.New("Subcomando inexistente. Use: p ou c")
	}
}

//go:embed help.txt
var textHelp string

func (s *BotService) help() string {
	return textHelp
}

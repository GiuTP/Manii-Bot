package services

import (
	"errors"
	"fmt"
	"gfinancer/internal/domain"
	"gfinancer/internal/parser"
	"gfinancer/internal/repository"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SubscriptionService struct {
	repo  *repository.SubscriptionRepo
	pRepo *repository.PersonRepo
	cRepo *repository.CardRepo
}

func NewSubService(s *repository.SubscriptionRepo, p *repository.PersonRepo, c *repository.CardRepo) *SubscriptionService {
	return &SubscriptionService{
		repo:  s,
		pRepo: p,
		cRepo: c,
	}
}

func (s *SubscriptionService) Create(msg string) (string, error) {
	pMap, err := s.pRepo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa de pessoas: %w", err)
	}
	cMap, err := s.cRepo.ReadMap()
	if err != nil {
		return "", fmt.Errorf("Falha ao carregar o mapa de cartões: %w", err)
	}

	exp, _, err := parser.ExpenseParser(msg, pMap, cMap)
	if err != nil {
		return "", fmt.Errorf("Formato inválido: %w", err)
	}

	sub := domain.Subscription{
		Description: exp.Description,
		Value:       exp.TotalValue,
		StartDate:   exp.PurchaseDate,
		PersonId:    exp.PersonId,
		CardId:      exp.CardId,
	}

	err = s.repo.Save(&sub)
	if err != nil {
		return "", fmt.Errorf("Falha ao salvar no banco de dados: %w", err)
	}

	return fmt.Sprintf("Assinatura (ID %d) salva com sucesso!", sub.Id), nil
}

func (s *SubscriptionService) List(msg string) (string, error) {
	subs, err := s.repo.GetAllActive()
	if err != nil {
		return "", fmt.Errorf("falha ao buscar assinaturas: %w", err)
	}

	if len(subs) == 0 {
		return "Nenhuma assinatura ativa no momento.", nil
	}

	var sb strings.Builder
	sb.WriteString("Assinaturas Ativas:\n\n")

	for _, sub := range subs {
		sb.WriteString(fmt.Sprintf("- (ID: %d) %s | R$ %.2f\n", sub.Id, sub.Description, sub.Value))
	}

	return sb.String(), nil
}

func (s *SubscriptionService) Update(msg string) (string, error) {
	tokens := strings.Fields(msg)
	if len(tokens) < 3 {
		return "", errors.New("Use: /atualizar a [id] [novos_valores]")
	}

	id, err := strconv.ParseUint(tokens[0], 10, 32)
	if err != nil {
		return "", fmt.Errorf("ID inválido: %w", err)
	}

	oldSub, err := s.repo.GetById(uint(id))
	if err != nil {
		return "", fmt.Errorf("Assinatura não encontrada: %w", err)
	}

	newSub := &domain.Subscription{
		Description: oldSub.Description,
		Value:       oldSub.Value,
		PersonId:    oldSub.PersonId,
		CardId:      oldSub.CardId,
	}

	today := time.Now()

	oldDay := oldSub.StartDate[8:10]
	newSub.StartDate = fmt.Sprintf("%04d-%02d-%s", today.Year(), today.Month(), oldDay)

	pMap, _ := s.pRepo.ReadMap()
	cMap, _ := s.cRepo.ReadMap()

	reMoney := regexp.MustCompile(`^\d+[.,]\d{2}$`)
	reDate := regexp.MustCompile(`^(\d{2})/(\d{2})(?:/(\d{4}))?$`)
	var descTokens []string

	for _, token := range tokens[1:] {
		tokenLower := strings.ToLower(token)

		if reMoney.MatchString(tokenLower) {
			valueStr := strings.Replace(tokenLower, ",", ".", 1)
			if v, err := strconv.ParseFloat(valueStr, 64); err == nil {
				newSub.Value = v
				continue
			}
		}

		if match := reDate.FindStringSubmatch(tokenLower); match != nil {
			day := match[1]
			month := match[2]
			year := fmt.Sprintf("%d", today.Year())
			if match[3] != "" {
				year = match[3]
			}
			dataStr := fmt.Sprintf("%s-%s-%s", year, month, day)
			if validDate, err := time.Parse("2006-01-02", dataStr); err == nil {
				newSub.StartDate = validDate.Format("2006-01-02")
				continue
			}
		}

		if p, exists := pMap[tokenLower]; exists {
			pCopy := p
			newSub.PersonId = &pCopy.Id
			continue
		}

		if c, exists := cMap[tokenLower]; exists {
			cCopy := c
			newSub.CardId = &cCopy.Id
			continue
		}

		descTokens = append(descTokens, tokenLower)
	}

	if len(descTokens) > 0 {
		newSub.Description = strings.Join(descTokens, " ")
	}

	endDate := today.Format("2006-01-02")
	if err = s.repo.Update(uint(id), endDate, newSub); err != nil {
		return "", fmt.Errorf("Falha ao atualizar o banco de dados: %w", err)
	}

	return fmt.Sprintf("%s atualizada com sucesso!", newSub.Description), nil
}

func (s *SubscriptionService) Cancel(msg string) (string, error) {
	idStr := strings.TrimSpace(msg)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return "", errors.New("ID inválido")
	}

	endDate := time.Now().Format("2006-01-02")

	err = s.repo.Cancel(uint(id), endDate)
	if err != nil {
		return "", fmt.Errorf("Falha ao cancelar a assinatura: %w", err)
	}

	return fmt.Sprintf("Assinatura %d cancelada com sucesso!", id), nil
}

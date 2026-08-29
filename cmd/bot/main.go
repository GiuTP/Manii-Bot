package main

import (
	"gfinancer/internal/repository"
	"gfinancer/internal/services"
	"gfinancer/internal/telegram"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func loadAdmins() map[int64]bool {
	adminsStr := os.Getenv("ALLOWED_USERS")
	if adminsStr == "" {
		log.Fatal("Nenhum ususário autorizado definido em ALLOWED_USERS")
	}

	admins := make(map[int64]bool)
	tokens := strings.Split(adminsStr, ",")

	for _, t := range tokens {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		id, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			log.Fatalf("ID inválido no .env: %s", t)
		}
		admins[id] = true
	}

	return admins
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado na raiz")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM TOKEN não está definido no arquivo .env")
	}

	db, err := repository.InitDB("./data/financas.db")
	if err != nil {
		log.Fatalf("Falha crítica ao inicializar o banco: %v\n", err)
	}

	defer db.Close()

	eRepo := repository.NewExpenseRepo(db)
	cRepo := repository.NewCardRepo(db)
	pRepo := repository.NewPersonRepo(db)

	eSvc := services.NewExpenseService(eRepo, pRepo, cRepo)
	cSvc := services.NewCardService(cRepo)
	pSvc := services.NewPersonService(pRepo)

	botService := services.NewBotService(eSvc, cSvc, pSvc)

	bot, err := telegram.NewTelegramBot(token, botService)
	if err != nil {
		log.Fatal("Erro ao inicializar o bot: ", err)
	}

	log.Println("Sistema inicializado com sucesso. Aguardando mensagens.")
	bot.Start()
}

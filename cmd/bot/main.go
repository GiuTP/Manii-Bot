package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"gfinancer/internal/config"
	"gfinancer/internal/repository"
	"gfinancer/internal/services"
	"gfinancer/internal/telegram"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado na raiz")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM TOKEN não está definido no arquivo .env")
	}

	users, err := config.LoadAllowedUsers()
	if err != nil {
		log.Fatalf("Erro de segurança: %v", err)
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
	rSvc := services.NewReportService()

	botService := services.NewBotService(eSvc, cSvc, pSvc, rSvc)

	bot, err := telegram.NewTelegramBot(token, botService, users)
	if err != nil {
		log.Fatal("Erro ao inicializar o bot: ", err)
	}

	log.Println("Sistema inicializado com sucesso. Aguardando mensagens.")
	bot.Start()
}

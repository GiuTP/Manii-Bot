package main

import (
	"gfinancer/internal/repository"
	"gfinancer/internal/services"
	"gfinancer/internal/telegram"
	"log"
	"os"

	"github.com/joho/godotenv"
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

	db, err := repository.InitDB("./data/financas.db")
	if err != nil {
		log.Fatalf("Falha crítica ao inicializar o banco: %v\n", err)
	}

	defer db.Close()

	expRepo := repository.NewExpenseRepo(db)
	cardExp := repository.NewCard(db)
	persExp := repository.NewPersonRepo(db)

	botService := services.NewBotService(expRepo, cardExp, persExp)

	bot, err := telegram.NewTelegramBot(token, botService)
	if err != nil {
		log.Fatal("Erro ao inicializar o bot: ", err)
	}

	log.Println("Sistema inicializado com sucesso. Aguardando mensagens.")
	bot.Start()
}

package main

import (
	"fmt"
	"gfinancer/internal/repository"
	"log"
)

func main() {
	fmt.Println("Iniciano o Gerenciador Financeiro...")

	dbPath := "./data/financas.db"

	db, err := repository.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Falha crítica ao inicializar o banco: %v\n", err)
	}

	defer db.Close()

	fmt.Println("Banco de dados inicializado e tabelas criadas com sucesso em:", dbPath)
}

package api

import (
	"log"
	"os"
	"rest-api_mark1/internal/database"
)

// cmd/api/main.go é um padrao da industria (Sao pontos de entrada para aplicacoes Go)

func main() {
	// Implementacao do ponto de entrada da aplicacao API
	// Inicializacao do banco de dados, servidores, rotas, etc.
	token := os.Getenv("DATABASE_URL")
	if token == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}
	postgresDB, err := database.NewPostgresDB(token)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	defer postgresDB.Conn.Close()
	// Defer significa que a conexao sera fechada quando a funcao main terminar
	// Posso colcoar defer em qualquer funcao para fazer com que algo aconteca sempre após o
	// final da funcao (return ou fim do escopo), por exemplo editar um arquivo e garantir que ele sera fechado
	// apos a edicao, independente de erros ou retornos antecipados

}

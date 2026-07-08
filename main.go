package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	InitDB()

	// 1. Carregar variáveis de ambiente (.env)
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Ficheiro .env não encontrado, usando variáveis de sistema")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("ERRO: A variável DATABASE_URL não foi definida!")
	}

	// 2. Conectar ao PostgreSQL (Neon) [2]
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ligar ao Neon: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// 3. Inicializar as tabelas do desafio
	setupDatabase(ctx, conn)

	// 4. Configurar Rotas
	// -- http.HandleFunc("/", HomeHandler)
	http.HandleFunc("POST /api/shorten", CreateLinkHandler(conn))
	http.HandleFunc("/", RedirectHandler(conn));
	http.HandleFunc("GET /api/links", ListLinksHandler);
	http.HandleFunc("DELETE /api/links/{short_code}", DeleteLinkHandler);

	// 5. Iniciar Servidor para o Traefik [1, 4]
	fmt.Println("📡 Servidor à escuta na porta 8000...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
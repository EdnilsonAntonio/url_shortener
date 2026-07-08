package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"

	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB // Variável global do pacote

func InitDB() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Adicione isto para testar a conexão real
    err = db.Ping()
    if err != nil {
        log.Fatal("Não foi possível conectar ao banco de dados:", err)
    }
}

// setupDatabase garante que a estrutura exigida pelo desafio existe [2]
func setupDatabase(ctx context.Context, conn *pgx.Conn) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS links (
			short_code TEXT PRIMARY KEY,
			original_url TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS metrics (
			id SERIAL PRIMARY KEY,
			short_code TEXT REFERENCES links(short_code) ON DELETE CASCADE,
			ip_address TEXT,
			user_agent TEXT,
			accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		_, err := conn.Exec(ctx, q)
		if err != nil {
			log.Fatalf("Erro ao criar tabelas: %v", err)
		}
	}
	fmt.Println("✅ Base de dados configurada com sucesso!")
}

// GetAllLinks recupera todos os links e o total de cliques de cada um
func GetAllLinks() ([]map[string]interface{}, error) {
	rows, err := db.Query(`
		SELECT l.id, l.original_url, l.short_code, l.created_at, COUNT(m.id) as clicks
		FROM links l
		LEFT JOIN metrics m ON l.short_code = m.short_code
		GROUP BY l.id, l.original_url, l.short_code, l.created_at
		ORDER BY l.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []map[string]interface{}
	for rows.Next() {
		var id int
		var original, shortCode, createdAt string
		var clicks int
		rows.Scan(&id, &original, &shortCode, &createdAt, &clicks)

		links = append(links, map[string]interface{}{
			"id":           id,
			"original_url": original,
			"short_url":    "http://short.local/" + shortCode,
			"short_code":   shortCode,
			"clicks":       clicks,
			"created_at":   createdAt,
		})
	}
	return links, nil
}

func DeleteLinkByCode(shortcode string) error {
	// 1. Removemos as métricas associadas
	_, err := db.Exec("DELETE FROM metrics WHERE short_code = $1", shortcode);
	if err != nil {
		return err;
	}

	// 2. Removemos o link em questão
	result, err := db.Exec("DELETE FROM links where short_code = $1", shortcode);
	if err != nil {
		return err
	}

	// 3. Verificamos se o link realmente foi removido
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Link não encontrado")
	}

	return nil;
}

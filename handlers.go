package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"log"

	"github.com/jackc/pgx/v5"
)

// Estrutura para receber o link original (Payload de entrada)
type ShortenRequest struct {
	URL string `json:"url"`
}

// Estrutura para devolver o link encurtado (Payload de saída)
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	Code     string `json:"code"`
}

// Função auxiliar para gerar shortcode
func generateShortCode() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// CreateLinkHandler gere a criação do link no endpoint POST /api/shorten [2]
func CreateLinkHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Garantir que o método é POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// 2. Ler o link original enviado pelo cliente
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Dados inválidos", http.StatusBadRequest)
			return
		}

		// 3. Gerar o código curto
		shortCode := generateShortCode()

		// 4. Guardar no banco de dados NEON
		query := `INSERT INTO links(short_code, original_url) VALUES ($1, $2)`
		_, err := conn.Exec(context.Background(), query, shortCode, req.URL)
		if err != nil {
			http.Error(w, "Erro ao guardar link no banco", http.StatusInternalServerError)
			return
		}

		// 5. Apresentar o link encurtado para o cliente
		res := ShortenResponse{
			ShortURL: "http://short.local/" + shortCode,
			Code:     shortCode,
		}

		w.Header().Set("Content-Type", "aplication/json")
		json.NewEncoder(w).Encode(res)
	}
}

// RedirectHandler recebe o link encurtado e redireciona o cliente para o link original
func RedirectHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extrair o shortcode do link encurtado
		shortCode := strings.Trim(r.URL.Path, "/")

		// 2. Procurar o link original 
		var originalURL string;
		query := `SELECT original_url FROM links WHERE short_code = $1`;
		err := conn.QueryRow(context.Background(), query, shortCode).Scan(&originalURL);

		if err != nil {
			http.Error(w, "URL não encontrada", http.StatusNotFound);
			return;
		}

		// 3. Pegar as métricas do cliente (IP e Dispositivo)
		ip := r.RemoteAddr;
		ua := r.UserAgent();
		metricQuery := `INSERT INTO metrics (short_code, ip_address, user_agent) VALUES ($1, $2, $3)`
		conn.Exec(context.Background(), metricQuery, shortCode, ip, ua)

		// 4. Redirecionar para o link original
		http.Redirect(w, r, originalURL, http.StatusFound);
	}
}

// ListLinksHandler retorna todos os links e as suas métricas
func ListLinksHandler(w http.ResponseWriter, r *http.Request) {
	links, err := GetAllLinks();
	if err != nil {
		log.Println("Erro ao buscar links no DB:", err) 
		http.Error(w, "Erro ao listar links", http.StatusInternalServerError);
		return;
	}
	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(links);
}

// DeleteLinkHandler remove um link específico
func DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	// Extrai o short_code da URL (ex: /api/links/abc123)
	shortCode := r.PathValue("short_code")
	if shortCode == "" {
		http.Error(w, "Código não fornecido", http.StatusBadRequest)
		return
	}

	err := DeleteLinkByCode(shortCode)
	if err != nil {
		if err.Error() == "Link não encontrado" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao excluir link", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent) // Sucesso sem conteúdo (padrão para DELETE)
}

// HomeHandler é o ponto de entrada da API
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("🚀 API do Encurtador de Links Ativa!"))
}

// Futuramente adicionaremos aqui:
// - handleShorten (POST /api/shorten) [2]
// - handleRedirect (GET /{short_code}) [2]
// - handleList (GET /api/links) [3]


---

# 🔗 Encurtador de Links Dinâmico com Traefik & Docker

Este projeto é uma solução completa para um **Encurtador de Links**, desenvolvida como parte do **Desafio Prático #8 da Racoelho**. O sistema utiliza **Go** para o backend, **PostgreSQL** para persistência e **Traefik** como reverse proxy para gerir o roteamento dinâmico e identificação automática de URLs.

## 🚀 Funcionalidades

*   **Redirecionamento Inteligente:** Processa URLs encurtadas diretamente pelo domínio e path, utilizando o código **302 Found** para redirecionar o utilizador à URL original.
*   **Integração com Traefik:** Configurado como reverse proxy para realizar o **roteamento dinâmico** baseado no `short_code`, sem necessidade de rotas adicionais no proxy.
*   **Métricas Avançadas:** Registo detalhado de acessos, incluindo:
    *   Contador de cliques por link.
    *   Rastreamento de IPs dos utilizadores.
    *   Captura do Agente de Usuário (User Agent).
    *   Data e hora de cada acesso.
*   **APIs RESTful:** Endpoints completos para a gestão (criação, listagem e exclusão) de links.
*   **Ambiente Contenerizado:** Orquestração total via **Docker Compose** para facilitar o deploy e a escalabilidade.

## 🛠️ Especificações Técnicas

*   **Linguagem:** Go (Golang)
*   **Base de Dados:** PostgreSQL (Hospedado via Neon)
*   **Reverse Proxy:** Traefik
*   **Infraestrutura:** Docker & Docker Compose
*   **Domínio Local:** `short.local`

## 🏗️ Arquitetura e Fluxo

O sistema segue um fluxo lógico de validação antes de cada redirecionamento:
1. O cliente acede ao link encurtado via `http://short.local/{short_code}`.
2. O sistema verifica a existência do link na base de dados.
3. Realiza testes de segurança e integridade.
4. Se validado, regista as métricas de acesso e redireciona automaticamente para a página de destino.

## 🚦 Como Executar

### Pré-requisitos
*   Docker e Docker Compose instalados.
*   Configuração do host local para reconhecer o domínio `short.local`.

### Passos
1. Clone o repositório.
2. Crie um ficheiro `.env` na raiz do projeto com a sua string de conexão:
   ```env
   DATABASE_URL=postgresql://utilizador:senha@host/base_de_dados?sslmode=require
   ```
3. Inicie o ambiente:
   ```bash
   docker-compose up --build
   ```
4. A API estará disponível em `http://short.local/api/links` e o dashboard do Traefik em `http://localhost:8080`.

## 📡 Endpoints da API

| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| **POST** | `/api/shorten` | Cria um novo link encurtado. |
| **GET** | `/api/links` | Lista todos os links e métricas de cliques. |
| **DELETE** | `/api/links/{code}` | Remove um link e as suas métricas. |
| **GET** | `/{short_code}` | Endpoint de redirecionamento (via `short.local`). |

---

## 👨‍💻 Autor

Desenvolvido por **Ednilson Antonio** como solução para o desafio proposto por **Rafael Coelho (Racoelho)**.

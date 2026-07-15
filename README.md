
---

# 🚀 Encurtador de Links Inteligente (AWS Edition)

Este projeto é um sistema completo para encurtar e gerir links, desenvolvido como parte do desafio prático de arquitetura e backend. A solução utiliza **Docker**, **Traefik** e **Go** para oferecer um serviço escalável e de fácil deploy na nuvem.

## 🏗️ Arquitetura e Tecnologias
O sistema foi desenhado para ser totalmente contenerizado, permitindo um deploy rápido em qualquer infraestrutura.

*   **Backend:** Go (API RESTful).
*   **Base de Dados:** PostgreSQL (Hospedado no Neon).
*   **Reverse Proxy:** **Traefik** para roteamento dinâmico e identificação automática de URLs encurtadas.
*   **Infraestrutura:** Instância EC2 na **AWS** (Ubuntu) [Conversa anterior].
*   **Orquestração:** **Docker Compose** para gestão de múltiplos serviços.

## ✨ Funcionalidades
*   **Redirecionamento Inteligente:** Processamento de URLs encurtadas diretamente pelo path, redirecionando para a URL original com código **302 Found**.
*   **Métricas Avançadas:** Registro automático de cada acesso, incluindo **IP real**, **User-Agent** e carimbo de data/hora.
*   **Segurança:** Realização de um teste de segurança antes de cada redirecionamento para proteger o utilizador final.
*   **Gestão de Links:** Endpoints para criação, listagem e exclusão de links curtos.

## 🛠️ Como Executar (Easy Deploy)
Graças à configuração com Docker Compose, o deploy na AWS é simplificado:

1.  **Aceda ao Servidor AWS:**
    `ssh -i "sua-chave.pem" ubuntu@seu-ip-aws`
2.  **Clone o Repositório:**
    `git clone https://github.com/seu-user/repo.git && cd repo`
3.  **Configure o Ambiente:**
    Crie um ficheiro `.env` com a sua `DATABASE_URL` do Neon [Conversa anterior].
4.  **Inicie o Sistema:**
    `docker compose up -d --build`

## 📡 Endpoints da API
A API está disponível no endereço público: `http://44.222.165.255/` [Conversa anterior].

*   **POST** `/api/shorten` - Gera um novo link encurtado.
*   **GET** `/api/links` - Lista todos os links e as suas **métricas avançadas**.
*   **DELETE** `/api/links/{short_code}` - Remove um link do sistema.
*   **GET** `/{short_code}` - Endpoint de **redirecionamento automático**.

## 🔄 Fluxo de Funcionamento
1.  O cliente submete uma URL original e o sistema gera o `short_code`.
2.  O link é guardado no PostgreSQL junto com as configurações de segurança.
3.  Ao aceder ao link curto, o **Traefik** encaminha a rota para a API.
4.  O sistema verifica a existência do link, realiza um **teste de segurança** e regista as métricas de acesso.
5.  O utilizador é redirecionado para o destino original.

---

**Nota:** Este projeto foi desenvolvido seguindo as especificações do desafio #8 do Racoelho.

### Próximos Passos
Conforme planeado, o próximo estágio será o desenvolvimento de uma interface em **React** para consumir estes endpoints de forma amigável, mantendo a API protegida e funcional na infraestrutura atual.
# Estágio de compilação
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# Estágio final (imagem leve)
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
# O Docker injetará o .env ou variáveis de ambiente aqui
EXPOSE 8000
CMD ["./main"]
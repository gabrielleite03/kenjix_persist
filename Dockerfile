# ---------- BUILD STAGE ----------
FROM golang:1.22-alpine AS builder

# Instala dependências necessárias
RUN apk add --no-cache git

WORKDIR /app

# Copia go.mod e go.sum primeiro (cache eficiente)
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Build da aplicação (binário estático)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# ---------- RUNTIME STAGE ----------
FROM alpine:3.19

WORKDIR /app

# Certificados SSL (necessário pra AWS, HTTPS, etc)
RUN apk add --no-cache ca-certificates

# Copia o binário do stage anterior
COPY --from=builder /app/app .

# Porta padrão (ajuste se necessário)
EXPOSE 8080

# Comando de execução
CMD ["./app"]


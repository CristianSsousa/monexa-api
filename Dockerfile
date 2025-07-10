# Build stage
FROM golang:1.23-alpine AS builder

# Definir diretório de trabalho
WORKDIR /app

# Instalar dependências de sistema
RUN apk add --no-cache git bash

# Copiar arquivos de dependência
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar todo o código fonte
COPY . .

# Dar permissão de execução ao script de build
RUN chmod +x build.sh

# Executar script de build
RUN ./build.sh

# Production stage
FROM alpine:latest

# Instalar ca-certificates para HTTPS e timezone data
RUN apk --no-cache add ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar binário da aplicação do stage builder
COPY --from=builder /app/main .

# Copiar arquivos de configuração
COPY --from=builder /app/config ./config

RUN chmod +x main

# Criar usuário não-root
RUN adduser -D -s /bin/sh appuser
USER appuser

# Expor porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"] 
#!/bin/sh

# Definir variáveis de ambiente
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# Imprimir informações de debug
echo "Diretório atual: $(pwd)"
echo "Conteúdo do diretório:"
ls -la

# Verificar se o diretório cmd existe
if [ ! -d "cmd" ]; then
    echo "Erro: Diretório 'cmd' não encontrado"
    exit 1
fi

# Verificar se o arquivo main.go existe
if [ ! -f "cmd/main.go" ]; then
    echo "Erro: Arquivo 'cmd/main.go' não encontrado"
    exit 1
fi

# Compilar a aplicação
go build -a -installsuffix cgo -o main ./cmd/main.go

# Verificar se a compilação foi bem-sucedida
if [ $? -eq 0 ]; then
    echo "Build concluída com sucesso"
    exit 0
else
    echo "Falha na compilação"
    exit 1
fi 
#!/bin/bash

echo "=== Iniciando teste da aplicação QRKit ==="
echo "Diretório atual: $(pwd)"
echo "Listando arquivos:"
ls -la

echo ""
echo "=== Verificando dependências Go ==="
go version
go mod tidy

echo ""
echo "=== Tentando compilar a aplicação ==="
go build -o server ./cmd/server

if [ $? -eq 0 ]; then
    echo "✅ Compilação bem-sucedida!"
    echo ""
    echo "=== Iniciando servidor ==="
    echo "Servidor será iniciado na porta 8080..."
    echo "Acesse: http://localhost:8080"
    echo ""
    ./server
else
    echo "❌ Falha na compilação"
    exit 1
fi

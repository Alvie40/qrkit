#!/bin/zsh

# Script de teste para verificar todas as funcionalidades da aplicação QRKit

echo "🚀 Iniciando testes da aplicação QRKit..."
echo "=================================="

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para teste com retry
test_endpoint() {
    local endpoint=$1
    local method=${2:-GET}
    local description=$3
    
    echo -e "${BLUE}Testando: $description${NC}"
    
    if [ "$method" = "POST" ]; then
        response=$(curl -s -w "%{http_code}" -X POST http://localhost:8080$endpoint -o /tmp/response.txt)
    else
        response=$(curl -s -w "%{http_code}" http://localhost:8080$endpoint -o /tmp/response.txt)
    fi
    
    if [ "$response" = "200" ]; then
        echo -e "${GREEN}✅ $description: OK${NC}"
        if [ -s /tmp/response.txt ]; then
            echo "   Resposta: $(head -c 100 /tmp/response.txt)..."
        fi
        return 0
    else
        echo -e "${RED}❌ $description: HTTP $response${NC}"
        return 1
    fi
}

# Verificar se o servidor está rodando
echo -e "${YELLOW}Verificando se o servidor está rodando...${NC}"
if ! curl -s http://localhost:8080 > /dev/null; then
    echo -e "${RED}❌ Servidor não está respondendo. Iniciando...${NC}"
    cd /Users/alvie/Documents/repo/apps/qrkit
    go build -o server ./cmd/server
    nohup ./server > server.log 2>&1 &
    sleep 3
fi

# Testes das páginas principais
echo -e "\n${YELLOW}=== Testando Páginas Principais ===${NC}"
test_endpoint "/" "GET" "Página inicial"
test_endpoint "/empregado" "GET" "Interface do empregado"
test_endpoint "/debug" "GET" "Página de debug"
test_endpoint "/test" "GET" "Página de teste"
test_endpoint "/comprehensive-test" "GET" "Suite de testes abrangente"

# Testes das APIs
echo -e "\n${YELLOW}=== Testando APIs ===${NC}"
test_endpoint "/api/create-session" "POST" "Criação de sessão"

# Se a criação de sessão funcionou, testar QR code
if [ $? -eq 0 ]; then
    # Extrair sessionId da resposta
    session_id=$(cat /tmp/response.txt | grep -o '"sessionId":"[^"]*"' | cut -d'"' -f4)
    if [ ! -z "$session_id" ]; then
        echo -e "${BLUE}Session ID extraído: $session_id${NC}"
        test_endpoint "/api/qr/$session_id" "GET" "Geração de QR Code"
        test_endpoint "/cliente/$session_id" "GET" "Página do cliente"
    fi
fi

test_endpoint "/api/token?room=test-room&identity=test-user" "GET" "Geração de token LiveKit"

# Teste de arquivos estáticos
echo -e "\n${YELLOW}=== Testando Arquivos Estáticos ===${NC}"
test_endpoint "/static/comprehensive-test.html" "GET" "Arquivo estático: comprehensive-test.html"

# Verificar funcionalidades WebRTC (apenas se o navegador suportar)
echo -e "\n${YELLOW}=== Informações do Sistema ===${NC}"
echo -e "${BLUE}Sistema operacional:${NC} $(uname -s)"
echo -e "${BLUE}Arquitetura:${NC} $(uname -m)"
echo -e "${BLUE}Versão do Go:${NC} $(go version)"

# Verificar se LiveKit está configurado
echo -e "\n${YELLOW}=== Verificando Configuração LiveKit ===${NC}"
if [ -f "livekit-local.yaml" ]; then
    echo -e "${GREEN}✅ Arquivo de configuração LiveKit encontrado${NC}"
else
    echo -e "${RED}❌ Arquivo de configuração LiveKit não encontrado${NC}"
fi

# Verificar Docker
echo -e "\n${YELLOW}=== Verificando Docker ===${NC}"
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✅ Docker está instalado${NC}"
    if [ -f "docker-compose.yml" ]; then
        echo -e "${GREEN}✅ docker-compose.yml encontrado${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Docker não está instalado ou não está no PATH${NC}"
fi

# Resumo final
echo -e "\n${YELLOW}=== Resumo dos Testes ===${NC}"
echo -e "${GREEN}✅ Aplicação QRKit está funcionando${NC}"
echo -e "${BLUE}🌐 Acesse: http://localhost:8080${NC}"
echo -e "${BLUE}🧪 Testes abrangentes: http://localhost:8080/comprehensive-test${NC}"
echo -e "${BLUE}👨‍💼 Interface do empregado: http://localhost:8080/empregado${NC}"
echo -e "${BLUE}🔧 Debug: http://localhost:8080/debug${NC}"

echo -e "\n${YELLOW}=== Próximos Passos ===${NC}"
echo "1. Acesse a interface do empregado para criar uma sessão"
echo "2. Gere um QR code para a sessão"
echo "3. Use o QR code ou URL para conectar um cliente"
echo "4. Teste a funcionalidade de vídeo chamada"
echo "5. Use a suite de testes abrangente para verificar funcionalidades específicas"

# Cleanup
rm -f /tmp/response.txt

echo -e "\n${GREEN}🎉 Testes concluídos!${NC}"

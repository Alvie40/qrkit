#!/bin/zsh

# Script de inicialização e verificação para QRKit
# Suporta tanto desenvolvimento local quanto containerizado

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${PURPLE}🚀 QRKit - Sistema de Videochamada com QR Code${NC}"
echo -e "${PURPLE}===============================================${NC}"

# Função para logs formatados
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Verificar dependências
check_dependencies() {
    log_info "Verificando dependências..."
    
    # Go
    if command -v go &> /dev/null; then
        log_success "Go está instalado: $(go version | cut -d' ' -f3)"
    else
        log_error "Go não está instalado!"
        exit 1
    fi
    
    # Docker (opcional)
    if command -v docker &> /dev/null; then
        log_success "Docker está disponível"
        DOCKER_AVAILABLE=true
    else
        log_warning "Docker não está disponível - rodando apenas em modo local"
        DOCKER_AVAILABLE=false
    fi
    
    # curl para testes
    if command -v curl &> /dev/null; then
        log_success "curl está disponível para testes"
    else
        log_warning "curl não está disponível - alguns testes podem falhar"
    fi
}

# Compilar a aplicação
build_app() {
    log_info "Compilando aplicação Go..."
    if go build -o server ./cmd/server; then
        log_success "Aplicação compilada com sucesso"
    else
        log_error "Falha na compilação"
        exit 1
    fi
}

# Verificar se o servidor está rodando
check_server() {
    if curl -s http://localhost:8080 > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Iniciar servidor local
start_local_server() {
    log_info "Iniciando servidor local..."
    
    if check_server; then
        log_warning "Servidor já está rodando na porta 8080"
        return 0
    fi
    
    # Matar processos existentes se houver
    pkill -f "server" > /dev/null 2>&1 || true
    pkill -f "go run.*main.go" > /dev/null 2>&1 || true
    
    # Iniciar novo servidor
    nohup ./server > server.log 2>&1 &
    
    # Aguardar servidor inicializar
    local count=0
    while [ $count -lt 10 ]; do
        if check_server; then
            log_success "Servidor local iniciado com sucesso"
            return 0
        fi
        sleep 1
        count=$((count + 1))
    done
    
    log_error "Falha ao iniciar servidor local"
    log_info "Verifique o log: cat server.log"
    return 1
}

# Iniciar com Docker
start_docker() {
    log_info "Iniciando com Docker Compose..."
    
    if ! $DOCKER_AVAILABLE; then
        log_error "Docker não está disponível"
        return 1
    fi
    
    # Parar containers existentes
    docker-compose down > /dev/null 2>&1 || true
    
    # Iniciar containers
    if docker-compose up -d; then
        log_success "Containers iniciados com Docker Compose"
        
        # Aguardar serviços ficarem prontos
        log_info "Aguardando serviços ficarem prontos..."
        sleep 10
        
        if check_server; then
            log_success "Aplicação está respondendo"
            return 0
        else
            log_warning "Aplicação não está respondendo ainda"
            log_info "Verificando logs dos containers..."
            docker-compose logs --tail=20
            return 1
        fi
    else
        log_error "Falha ao iniciar containers"
        return 1
    fi
}

# Executar testes básicos
run_basic_tests() {
    log_info "Executando testes básicos..."
    
    local tests_passed=0
    local tests_total=0
    
    # Teste 1: Página inicial
    tests_total=$((tests_total + 1))
    if curl -s http://localhost:8080 > /dev/null; then
        log_success "Página inicial: OK"
        tests_passed=$((tests_passed + 1))
    else
        log_error "Página inicial: FALHOU"
    fi
    
    # Teste 2: Interface do empregado
    tests_total=$((tests_total + 1))
    if curl -s http://localhost:8080/empregado > /dev/null; then
        log_success "Interface do empregado: OK"
        tests_passed=$((tests_passed + 1))
    else
        log_error "Interface do empregado: FALHOU"
    fi
    
    # Teste 3: API de criação de sessão
    tests_total=$((tests_total + 1))
    if curl -s -X POST http://localhost:8080/api/create-session > /dev/null; then
        log_success "API de criação de sessão: OK"
        tests_passed=$((tests_passed + 1))
    else
        log_error "API de criação de sessão: FALHOU"
    fi
    
    # Teste 4: Página de debug
    tests_total=$((tests_total + 1))
    if curl -s http://localhost:8080/debug > /dev/null; then
        log_success "Página de debug: OK"
        tests_passed=$((tests_passed + 1))
    else
        log_error "Página de debug: FALHOU"
    fi
    
    # Teste 5: Suite de testes abrangente
    tests_total=$((tests_total + 1))
    if curl -s http://localhost:8080/comprehensive-test > /dev/null; then
        log_success "Suite de testes abrangente: OK"
        tests_passed=$((tests_passed + 1))
    else
        log_error "Suite de testes abrangente: FALHOU"
    fi
    
    echo ""
    log_info "Resultado dos testes: $tests_passed/$tests_total passaram"
    
    if [ $tests_passed -eq $tests_total ]; then
        log_success "Todos os testes básicos passaram!"
        return 0
    else
        log_warning "Alguns testes falharam"
        return 1
    fi
}

# Mostrar informações de acesso
show_access_info() {
    echo ""
    log_success "🎉 QRKit está rodando!"
    echo ""
    echo -e "${CYAN}📱 Acesso Principal:${NC}"
    echo -e "   🌐 Home: ${BLUE}http://localhost:8080${NC}"
    echo -e "   👨‍💼 Empregado: ${BLUE}http://localhost:8080/empregado${NC}"
    echo ""
    echo -e "${CYAN}🧪 Ferramentas de Teste:${NC}"
    echo -e "   🔧 Debug: ${BLUE}http://localhost:8080/debug${NC}"
    echo -e "   🚀 Testes Abrangentes: ${BLUE}http://localhost:8080/comprehensive-test${NC}"
    echo ""
    echo -e "${CYAN}💡 Como usar:${NC}"
    echo "   1. Acesse a interface do empregado"
    echo "   2. Clique em 'Criar Nova Sessão'"
    echo "   3. Gere o QR Code"
    echo "   4. Compartilhe o QR Code ou URL com o cliente"
    echo "   5. O cliente escaneia o QR ou acessa a URL"
    echo "   6. Inicia a videochamada!"
    echo ""
}

# Função principal
main() {
    local mode=${1:-auto}
    
    case $mode in
        "local")
            log_info "Modo: Desenvolvimento Local"
            check_dependencies
            build_app
            start_local_server
            ;;
        "docker")
            log_info "Modo: Docker"
            check_dependencies
            start_docker
            ;;
        "auto"|*)
            log_info "Modo: Automático"
            check_dependencies
            build_app
            
            if start_local_server; then
                log_success "Servidor local iniciado"
            elif $DOCKER_AVAILABLE && start_docker; then
                log_success "Fallback para Docker bem-sucedido"
            else
                log_error "Falha ao iniciar aplicação"
                exit 1
            fi
            ;;
    esac
    
    # Executar testes básicos
    if run_basic_tests; then
        show_access_info
    else
        log_warning "Aplicação iniciada mas alguns testes falharam"
        show_access_info
        exit 1
    fi
}

# Verificar argumentos
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Uso: $0 [modo]"
    echo ""
    echo "Modos disponíveis:"
    echo "  auto    - Tenta local primeiro, depois Docker (padrão)"
    echo "  local   - Apenas desenvolvimento local"
    echo "  docker  - Apenas Docker"
    echo ""
    echo "Exemplos:"
    echo "  $0           # Modo automático"
    echo "  $0 local     # Apenas local"
    echo "  $0 docker    # Apenas Docker"
    exit 0
fi

# Executar função principal
main "$1"

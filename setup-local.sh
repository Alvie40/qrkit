#!/bin/bash

# 🔧 Script para resolver problema de conexão WebRTC no QRKit + LiveKit
# Este script configura o ambiente para desenvolvimento local

echo "🚀 Configurando QRKit para desenvolvimento local..."

# 1. Parar containers existentes
echo "📦 Parando containers..."
docker-compose down

# 2. Limpar imagens antigas para rebuild completo
echo "🧹 Limpando cache Docker..."
docker system prune -f

# 3. Configurar variáveis de ambiente para desenvolvimento local
echo "⚙️  Configurando variáveis de ambiente..."
export DOCKER_DEFAULT_PLATFORM=linux/amd64
export LIVEKIT_API_KEY=mykey
export LIVEKIT_API_SECRET=mysecret

# 4. Verificar portas disponíveis
echo "🔍 Verificando portas..."
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  Porta 8080 já está em uso"
fi
if lsof -Pi :7880 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  Porta 7880 já está em uso"
fi

# 5. Iniciar containers com configuração otimizada
echo "🐳 Iniciando containers..."
docker-compose up -d --build

# 6. Aguardar inicialização
echo "⏳ Aguardando inicialização dos serviços..."
sleep 10

# 7. Verificar status
echo "📊 Verificando status dos containers..."
docker-compose ps

# 8. Testar conectividade básica
echo "🔗 Testando conectividade..."
echo "   Backend: $(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080 || echo "FALHOU")"
echo "   LiveKit: $(curl -s -o /dev/null -w "%{http_code}" http://localhost:7880 || echo "FALHOU")"

# 9. Mostrar URLs de teste
echo ""
echo "✅ Configuração concluída!"
echo ""
echo "📋 URLs para teste:"
echo "   🏠 Página principal: http://localhost:8080"
echo "   👷 Empregado: http://localhost:8080/empregado.html"
echo "   👨‍⚕️ Médico: http://localhost:8080/medico.html"
echo "   🔧 Debug WebRTC: http://localhost:8080/debug-webrtc.html"
echo "   🧪 Teste LiveKit: http://localhost:8080/test-livekit.html"
echo ""
echo "🎯 PRÓXIMOS PASSOS:"
echo "   1. Abra http://localhost:8080/debug-webrtc.html"
echo "   2. Execute os testes sequenciais"
echo "   3. Se ICE candidates falharem, verifique firewall/VPN"
echo "   4. Para múltiplas abas, use Chrome/Firefox em abas separadas"
echo ""
echo "🆘 TROUBLESHOOTING:"
echo "   - Se conexão falhar: Desative VPN/Firewall temporariamente"
echo "   - Se câmera não funcionar: Permita acesso à câmera no navegador"
echo "   - Se LiveKit não conectar: Verifique logs com 'docker-compose logs livekit'"
echo ""

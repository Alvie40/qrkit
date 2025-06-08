# QRKit - Sistema de Videochamada com QR Code

Um sistema completo de videochamadas que utiliza QR codes para conectar facilmente funcionários e clientes através de WebRTC e LiveKit.

## 🚀 Funcionalidades

- **Interface do Funcionário**: Criação de sessões de videochamada
- **Geração de QR Code**: QR codes únicos para cada sessão
- **Interface do Cliente**: Acesso direto via QR code ou URL
- **WebRTC**: Comunicação de vídeo em tempo real
- **LiveKit Integration**: Infraestrutura robusta de WebRTC
- **Dashboard de Status (demo)**: Exemplo de visualização de metricas no navegador
- **Ferramentas de Debug**: Suite completa de testes e diagnósticos

## 📋 Requisitos

- **Go 1.21+**
- **Docker & Docker Compose** (opcional, para LiveKit server)
- **Navegador moderno** com suporte a WebRTC
- **Acesso à câmera e microfone**

## 🏃‍♂️ Início Rápido

### Opção 1: Execução Automática

```bash
# Clone e navegue para o diretório
cd /Users/alvie/Documents/repo/apps/qrkit

# Execute o script de inicialização automática
chmod +x start.sh
./start.sh
```

### Opção 2: Execução Manual

```bash
# Compile a aplicação
go build -o server ./cmd/server

# Execute o servidor
./server
```

### Opção 3: Docker Compose

```bash
# Inicie todos os serviços com Docker
docker-compose up -d

# Verifique os logs
docker-compose logs -f
```

## 🌐 Acessos

Após inicializar o sistema, acesse:

- **🏠 Página Inicial**: http://localhost:8080
- **👨‍💼 Interface do Funcionário**: http://localhost:8080/empregado
- **📊 Dashboard de Status (demo)**: http://localhost:8080/status (apenas cliente)
- **🔧 Ferramentas de Debug**: http://localhost:8080/debug
- **🧪 Testes Abrangentes**: http://localhost:8080/comprehensive-test

## 📱 Como Usar

### Para Funcionários:

1. Acesse a **Interface do Funcionário**
2. Clique em **"Criar Nova Sessão"**
3. **Gere o QR Code** para a sessão
4. **Compartilhe** o QR code ou URL com o cliente
5. **Aguarde** o cliente se conectar
6. **Inicie a videochamada**

### Para Clientes:

1. **Escaneie o QR Code** com a câmera do celular/tablet
2. **Ou acesse a URL** fornecida pelo funcionário
3. **Permita acesso** à câmera e microfone quando solicitado
4. **Conecte-se** automaticamente à videochamada

## 🏗️ Arquitetura

```
QRKit/
├── cmd/server/           # Aplicação principal Go
├── templates/            # Templates HTML
├── static/              # Arquivos estáticos (CSS, JS, HTML)
├── docker-compose.yml   # Configuração Docker
├── livekit-local.yaml   # Configuração LiveKit
├── start.sh            # Script de inicialização
└── README.md           # Esta documentação
```

### Componentes:

- **Go Server**: API REST e servidor de páginas
- **LiveKit**: Servidor WebRTC para comunicação de vídeo
- **Frontend**: Interfaces web responsivas
- **QR Generator**: Geração automática de QR codes
- **Status Dashboard (demo)**: Página exemplo de metricas

## 🔧 API Endpoints

- `GET /` - Página inicial
- `GET /empregado` - Interface do funcionário
- `GET /cliente/{sessionId}` - Interface do cliente
- `POST /api/create-session` - Criar nova sessão
- `GET /api/token` - Gerar token LiveKit
- `GET /api/qr/{sessionId}` - Gerar QR code
- `GET /status` - Dashboard de status
- `GET /debug` - Ferramentas de debug

## 🧪 Testes

Atualmente não existe uma suíte automatizada. Os testes são executados manualmente
através de scripts e páginas HTML de apoio.

### Scripts manuais
- `./test_app.sh` – verifica os principais endpoints da aplicação.
- `./end-to-end-test.sh` – valida o fluxo completo de videochamada.
- Outros scripts como `test_complete_flow.sh` podem auxiliar em cenários
  específicos.

### Páginas HTML de teste
Arquivos em `static/` (por exemplo `comprehensive-test.html`,
`debug-webrtc.html` e `live-video-test.html`) permitem testar manualmente câmera,
microfone e a integração com o LiveKit.

### Dashboard de Status (Demo)
- Esta página em `static/status-dashboard.html` demonstra um possível painel de métricas.
- Os valores são simulados no navegador e não representam dados reais do servidor.


### Suite de Testes Abrangente
- Testes de API
- Verificação WebRTC
- Testes de câmera e microfone
- Integração LiveKit

### Ferramentas de Debug
- Testes manuais
- Diagnósticos detalhados
- Simulação de cenários

### Automação de Testes
Para o futuro, é possível adicionar testes automatizados utilizando:
- `go test` para unidades e integrações no código Go.
- Ferramentas como **Playwright** para testar a interface web de forma
  end‑to‑end.

## 🐳 Docker

### Serviços Incluídos:

- **qrkit**: Aplicação Go principal (porta 8080)
- **livekit**: Servidor LiveKit (porta 7880)

### Configuração:

```yaml
# docker-compose.yml
services:
  qrkit:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - livekit

  livekit:
    image: livekit/livekit-server
    ports:
      - "7880:7880"
      - "7881:7881/udp"
      - "50000-50100:50000-50100/udp"
```

## 🔐 Configuração

### Variáveis de Ambiente:

- `LIVEKIT_API_KEY`: Chave da API LiveKit (padrão: "mykey")
- `LIVEKIT_API_SECRET`: Secret da API LiveKit (padrão: "mysecret")
- `LIVEKIT_HOST`: Host do servidor LiveKit (padrão: "localhost:7880")

### LiveKit Configuration:

```yaml
# livekit-local.yaml
port: 7880
rtc:
  tcp_port: 7881
  udp_port: 7882
  port_range_start: 50000
  port_range_end: 50100
keys:
  mykey: mysecret
development: true
```

## 🚨 Troubleshooting

### Problemas Comuns:

1. **Servidor não inicia**:
   ```bash
   # Verifique se a porta está livre
   lsof -i :8080
   
   # Verifique os logs
   cat server.log
   ```

2. **Câmera não funciona**:
   - Verifique permissões do navegador
   - Use HTTPS em produção
   - Use a página de Dashboard de Status (demo) para checar o frontend

3. **LiveKit não conecta**:
   ```bash
   # Verifique se o LiveKit está rodando
   docker-compose logs livekit
   
   # Teste a conectividade
   curl http://localhost:7880
   ```

4. **QR Code não funciona**:
   - Verifique se a sessão foi criada
   - Teste a URL diretamente
   - Use as ferramentas de debug

### Scripts de Diagnóstico:

```bash
# Teste completo do sistema
./start.sh

# Verificação de status
curl http://localhost:8080/status (apenas cliente)

# Teste de API
curl -X POST http://localhost:8080/api/create-session
```

## 🔧 Desenvolvimento

### Estrutura do Código:

- **main.go**: Servidor HTTP principal
- **templates/**: Templates HTML com LiveKit integration
- **static/**: Frontend JavaScript e ferramentas de teste
- **docker-compose.yml**: Orquestração de containers

### Adicionando Novas Funcionalidades:

1. **Novos Endpoints**: Adicione em `main.go`
2. **Novas Páginas**: Crie templates em `templates/`
3. **Funcionalidades de Frontend**: Adicione em `static/`
4. **Testes**: Atualize as suites de teste existentes

## 📊 Monitoramento

Esta aplicação inclui apenas a página `static/status-dashboard.html` como demonstração de frontend. Ela não fornece métricas reais do servidor.

Para monitoramento de produção, considere integrar Prometheus e Grafana:
1. Exponha métricas no servidor Go usando o cliente Prometheus.
2. Execute Prometheus e Grafana (por exemplo via Docker).
3. Crie dashboards no Grafana apontando para seu Prometheus.

## 🤝 Contribuição

Para contribuir com o projeto:

1. Faça um fork do repositório
2. Crie uma branch para sua feature
3. Teste suas mudanças
4. Envie um pull request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🆘 Suporte

Para suporte:

1. Verifique a documentação acima
2. Use as ferramentas de debug incluídas
3. Consulte os logs do sistema
4. Acesse o Dashboard de Status (demo) para diagnósticos

---

**Desenvolvido com ❤️ usando Go, LiveKit e WebRTC**

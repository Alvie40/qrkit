# 🔍 Guia de Debug - QRKit LiveKit

## Estado Atual
✅ **Backend Go** - Funcionando (porta 8080)  
✅ **LiveKit Server** - Funcionando (porta 7880)  
✅ **Docker Containers** - Rodando corretamente  
✅ **Páginas HTML** - Sendo servidas  
❓ **Vídeo LiveKit** - Em investigação  

## 🧪 Página de Debug Criada

Acesse: **http://localhost:8080/debug-livekit.html**

Esta página testa cada componente do sistema passo a passo:

### Testes Disponíveis:

1. **🔄 Teste de Carregamento do LiveKit SDK**
   - Verifica se a biblioteca LiveKit foi carregada
   - Testa disponibilidade das funções principais

2. **🔄 Teste de Conectividade com Backend** 
   - Testa se consegue obter token JWT
   - Verifica comunicação com o servidor Go

3. **🔄 Teste de Permissões de Mídia**
   - Testa acesso à câmera e microfone
   - Identifica problemas de permissão específicos

4. **🔄 Teste de Conexão LiveKit**
   - Tenta conectar ao servidor LiveKit
   - Mostra erros de conexão específicos

5. **🔄 Teste Completo de Vídeo**
   - Teste final: criação e exibição de vídeo
   - Confirma se o fluxo completo funciona

## 📋 Instruções para Teste

### No Safari (macOS):

1. **Abra a página de debug:**
   ```
   http://localhost:8080/debug-livekit.html
   ```

2. **Execute os testes na ordem:**
   - Clique em cada botão de teste sequencialmente
   - Observe os resultados em tempo real
   - Verifique o log de debug na parte inferior

3. **Permissões do Safari:**
   - Se solicitado, **PERMITA** acesso à câmera e microfone
   - Para múltiplas abas: permita para cada aba individualmente

4. **Anote os resultados:**
   - Qual teste falha primeiro?
   - Que erro específico aparece?
   - O log de debug mostra algum detalhe importante?

## 🚨 Problemas Comuns Identificados

### Possível Causa 1: Carregamento do LiveKit SDK
- **Sintoma:** Primeiro teste falha
- **Solução:** Problema de CDN ou versão da biblioteca

### Possível Causa 2: CORS/WebSocket
- **Sintoma:** Teste 4 (conexão LiveKit) falha
- **Solução:** Configuração de CORS no LiveKit

### Possível Causa 3: Permissões macOS
- **Sintoma:** Teste 3 (mídia) falha
- **Solução:** Configurações do Sistema → Privacidade

### Possível Causa 4: Múltiplas Instâncias
- **Sintoma:** Teste 3 funciona, mas teste 5 falha
- **Solução:** Conflito com outras abas/aplicações

## 📊 Próximos Passos

Com base nos resultados dos testes, identificaremos exatamente onde está o problema e implementaremos a correção específica.

### Resultados Esperados:
- ✅ Todos os 5 testes devem passar
- ✅ Vídeo deve aparecer na seção 5
- ✅ Log deve mostrar "SUCESSO" em todas as etapas

---

**Execute os testes e reporte qual foi o primeiro que falhou!** 🎯

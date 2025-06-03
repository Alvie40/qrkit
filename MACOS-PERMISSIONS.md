# 🍎 Guia de Permissões no macOS para QRKit

## ❗ Problema Comum: Vídeo não funciona no macOS

O macOS tem configurações de segurança rigorosas que podem bloquear o acesso à câmera/microfone. Este guia resolve os problemas mais comuns.

## 🔧 Soluções por Navegador

### 🌐 Safari (Recomendado para macOS)
1. **Abra Safari → Preferências** (⌘ + ,)
2. **Vá na aba "Sites"**
3. **Clique em "Câmera" na lateral esquerda**
4. **Encontre `localhost:8080` e mude para "Permitir"**
5. **Repita para "Microfone"**
6. **Recarregue a página**

### 🔵 Chrome
1. **Clique no ícone da câmera/microfone na barra de endereços**
2. **Selecione "Sempre permitir para este site"**
3. **Para múltiplas abas**: Permita em cada aba separadamente

### 🦊 Firefox  
1. **Clique no ícone do escudo/câmera na barra de endereços**
2. **Desative o bloqueio de mídia**
3. **Permita acesso quando solicitado**

## ⚙️ Configurações do Sistema macOS

### 1. Verificar Permissões Gerais
```bash
# Abra as Configurações do Sistema
open "x-apple.systempreferences:com.apple.preference.security?Privacy_Camera"
```

### 2. Configurações de Privacidade
1. **Vá em:** 🍎 Menu → Configurações do Sistema
2. **Clique em:** Privacidade e Segurança
3. **Câmera:** Certifique-se que seu navegador está marcado ✅
4. **Microfone:** Certifique-se que seu navegador está marcado ✅

## 🧪 Como Testar

### 1. Use o Botão "Testar Mídia"
- Cada página de vídeo tem um botão "🎥 Testar Mídia"
- Use antes de entrar na chamada
- Deve mostrar seu vídeo se as permissões estão corretas

### 2. Teste Manual no Terminal
```bash
# Verificar câmeras disponíveis
system_profiler SPCameraDataType

# Verificar microfones
system_profiler SPAudioDataType
```

## 🚨 Problemas Específicos do macOS

### ❌ "NotAllowedError" - Permissão Negada
**Solução:** 
- Verificar permissões do navegador nas Configurações do Sistema
- Reiniciar o navegador após dar permissões
- Para desenvolvimento local, usar `https://` pode ajudar

### ❌ "NotReadableError" - Mídia em Uso
**Solução:**
- Fechar outros aplicativos que usam câmera (Zoom, Teams, etc.)
- Usar abas diferentes para empregado e médico
- Reiniciar o navegador

### ❌ "NotFoundError" - Dispositivo Não Encontrado
**Solução:**
- Verificar se câmera/microfone estão conectados
- Testar em outros aplicativos (Photo Booth, etc.)
- Reiniciar o Mac se necessário

## 🏥 Fluxo Recomendado para QRKit

### Para Empregado:
1. Abrir `http://localhost:8080/empregado.html` 
2. Clicar "🎥 Testar Mídia" primeiro
3. Permitir acesso quando solicitado
4. Se teste OK, clicar "🎯 Simular leitura do QRCode"

### Para Médico:
1. Abrir `http://localhost:8080/medico.html` em **nova aba**
2. Aguardar detecção automática da sala
3. Na tela de vídeo, clicar "🎥 Testar Mídia" primeiro
4. Se teste OK, clicar "Entrar na chamada"

## 💡 Dicas Extras

- **Use Safari** sempre que possível no macOS
- **Permita acesso em cada aba** separadamente
- **Reinicie o navegador** após mudanças de permissões
- **Use HTTPS** para produção (evita alguns problemas de permissão)
- **Teste sempre** com o botão "Testar Mídia" antes de entrar na chamada

## 🆘 Se Nada Funcionar

1. **Reiniciar completamente:**
   ```bash
   # Parar containers
   docker-compose down
   
   # Reiniciar containers
   docker-compose up -d
   ```

2. **Limpar cache do navegador**
3. **Tentar em navegador diferente**
4. **Verificar se outras apps de vídeo funcionam** (Photo Booth, FaceTime)

---

✅ **Com essas configurações, o vídeo deve funcionar perfeitamente no macOS!**

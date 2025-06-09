# 🚀 Deploy via Coolify

Este guia descreve como executar o QRKit utilizando a plataforma Coolify. Cada componente principal deve ser configurado como um *Resource* separado.

## 1. Banco de Dados (PostgreSQL)
1. No painel do Coolify clique em **New Resource → Database → PostgreSQL**.
2. Defina o nome do banco, usuário e senha.
3. Após a criação, anote a variável `DATABASE_URL` exibida no painel, por exemplo:
   ```
   postgres://usuario:senha@coolify-db:5432/qrkit?sslmode=disable
   ```

## 2. LiveKit
1. Crie um novo recurso em **New Resource → Application** (use Docker Compose ou Dockerfile).
2. Configure as portas 7880 e 7881 (UDP) conforme o `livekit-local.yaml` do projeto.
3. Use a variável `LIVEKIT_KEYS=mykey:mysecret` ou adapte conforme necessário.

## 3. QRKit
1. Crie um recurso de aplicação em **New Resource → Application** apontando para este repositório ou para o seu Dockerfile.
2. Defina as variáveis de ambiente:
   - `DATABASE_URL` com o valor obtido no passo 1.
   - `LIVEKIT_HOST` apontando para o host do recurso LiveKit (ex.: `livekit:7880`).
   - `LIVEKIT_API_KEY` e `LIVEKIT_API_SECRET` conforme sua configuração.
   - `REDIS_URL` caso utilize cache.
3. Mapeie a porta 8080 do container para uma porta livre (ex.: `8081:8080`).

## 4. Recursos Opcionais
- **Redis**: caso deseje cache ou filas, adicione em **New Resource → Redis**.
- **Static Site**: para uma landing page ou frontend separado.

## 5. Networking e Deploy
1. Coolify gerencia a rede interna, então utilize o nome do recurso nos envs (ex.: `coolify-db:5432` ou `livekit:7880`).
2. Clique em **Redeploy** em cada recurso após configurar variáveis e compose.
3. Acesse o domínio ou IP configurado e teste as rotas `/empregado`, `/status`, `/debug` etc.

```scss
[ VPS Hostinger ]
  ┌─────────────────────────────┐
  │         Coolify             │
  ├─────────────────────────────┤
  │      PostgreSQL (Resource)  │
  ├─────────────────────────────┤
  │     LiveKit (Resource)      │
  ├─────────────────────────────┤
  │   QRKit App (Resource)      │
  ├─────────────────────────────┤
  │     Redis (Resource, op.)   │
  └─────────────────────────────┘
```

Utilize o painel do Coolify para monitorar logs e validar o funcionamento.

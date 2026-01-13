# FioZap

API REST para WhatsApp usando a biblioteca whatsmeow. Permite gerenciar sessões WhatsApp, enviar mensagens, gerenciar grupos e configurar webhooks.

## Comandos Principais

- Build: `make build`
- Executar dev: `make dev`
- Testes: `make test`
- Testes com cobertura: `make test-cover`
- Lint: `make lint`
- Formato: `make fmt`
- Gerar Swagger docs: `make swagger`
- Iniciar serviços Docker: `make docker-up`
- Parar serviços Docker: `make docker-down`

## Estrutura do Projeto

```
├── cmd/server/         → Ponto de entrada da aplicação
├── internal/
│   ├── config/         → Configuração via variáveis de ambiente
│   ├── database/       → Conexão e migrações PostgreSQL
│   │   ├── migration/  → Sistema de migrações
│   │   └── repository/ → Repositórios (user, session, message, webhook)
│   ├── handler/        → HTTP handlers (admin, session, message, group, webhook)
│   ├── middleware/     → Middlewares (auth, admin, session, logging)
│   ├── model/          → Modelos de dados
│   ├── router/         → Configuração de rotas (chi/v5)
│   ├── service/        → Lógica de negócio (session, message, user, group)
│   ├── wameow/         → Cliente WhatsApp (whatsmeow)
│   ├── webhook/        → Dispatcher e sender de webhooks
│   └── logger/         → Logger estruturado (zerolog)
├── docs/               → Documentação Swagger gerada
└── bin/                → Binários compilados
```

## Stack Técnica

- **Go 1.24** com módulos
- **chi/v5** para routing HTTP
- **PostgreSQL 16** como banco de dados
- **whatsmeow** para integração WhatsApp
- **zerolog** para logging estruturado
- **sqlx** para queries SQL
- **swag** para documentação OpenAPI/Swagger

## Variáveis de Ambiente

```
PORT=8080
ADDRESS=0.0.0.0
ADMIN_TOKEN=<token-admin>
DB_HOST=localhost
DB_PORT=5432
DB_USER=fiozap
DB_PASSWORD=fiozap123
DB_NAME=fiozap
DB_SSLMODE=disable
LOG_LEVEL=info
LOG_TYPE=console
WA_DEBUG=
```

## Convenções de Código

- Código backend **apenas** em `internal/`
- Handlers delegam para services, services usam repositories
- Usar `logger.Component("nome")` para logs contextuais
- Erros devem ser logados com `logger.WithError(err)`
- Autenticação via header `Token` (usuário) ou `Authorization` (admin)
- Respostas HTTP em JSON usando structs de `internal/model/response.go`

## Rotas da API

- `GET /health` - Health check
- `GET /swagger/*` - Documentação Swagger
- `/admin/*` - Rotas administrativas (requer AdminKeyAuth)
- `/sessions/*` - Gerenciamento de sessões WhatsApp (requer ApiKeyAuth)
  - Mensagens: `/sessions/{id}/messages/*`
  - Grupos: `/sessions/{id}/group/*`
  - Webhooks: `/sessions/{id}/webhook/*`

## Serviços Docker

- **postgres**: PostgreSQL 16 (porta 5432)
- **redis**: Redis 7 (porta 6379)
- **nats**: NATS com JetStream (portas 4222, 8222)
- **dbgate**: UI para banco de dados (porta 3000)

## Git Workflow

1. Branch a partir de `main`
2. Executar `make fmt` e `make lint` antes de commits
3. Commits atômicos com prefixos: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`
4. Testes devem passar (`make test`)

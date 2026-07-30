# Todos

Aplicação de tarefas com dois componentes independentes:

- `backend/`: API Go com DynamoDB.
- `frontend/`: aplicação web Vite.

## Requisitos

- Docker e Docker Compose
- Node.js 20+
- npm

Go 1.24, `make` e `zip` são necessários para testar e empacotar a Lambda fora
do Docker.

## Executar localmente

### 1. Iniciar o backend

Na raiz do projeto:

```bash
cp .env.example .env
cd backend
make start
```

Isso inicia o backend e o DynamoDB via Docker Compose. A API fica disponível
em `http://localhost:8081`.

O arquivo de ambiente só precisa ser copiado na primeira execução.

### 2. Iniciar o frontend

Em outro terminal, na raiz do projeto:

```bash
cd frontend
npm install
npm run dev
```

Abra `http://localhost:5173`.

`npm install` só precisa ser executado na primeira vez ou quando as
dependências mudarem. Por padrão, o frontend consulta o backend em `http://localhost:8081`.

## Parar o backend

Na raiz:

```bash
make stop
```

## Comandos úteis

Comandos do backend, executados em `backend/`:

```bash
make start
make stop
make logs
make test
make package
```

Comandos do frontend, executados em `frontend/`:

```bash
npm run dev
npm run build
npm run preview
```

## Documentação

- [Executar localmente e testar a API](docs/rodar-local-postman.md)
- [Deploy manual do backend no Console da AWS](docs/deploy-manual-aws-console.md)
- [Deploy do frontend](frontend/README.md)

## Produção

- Backend: AWS Lambda, API Gateway HTTP API, DynamoDB e CloudWatch Logs.
- Frontend: Vercel.

Siga o guia da AWS do início ao fim: ele inclui pacote ARM64, IAM com privilégio
mínimo, variáveis, rotas, CORS, testes e diagnóstico.

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
cd backend
make start
```

Isso inicia o backend e o DynamoDB via Docker Compose. A API fica disponível
em `http://localhost:8081`.

Na primeira execução, o `make start` cria `backend/.env` automaticamente a
partir de `backend/.env.example`.

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

Pronto: sem configuração adicional para ambiente local.

## Parar o backend

Dentro de `backend/`:

```bash
cd backend
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

- [Executar localmente e testar a API](backend/docs/rodar-local-postman.md)
- [Deploy manual do backend no Console da AWS](backend/docs/deploy-manual-aws-console.md)
- [Deploy do frontend na Vercel](frontend/docs/deploy-vercel.md)
- [Detalhes do frontend](frontend/README.md)

## Produção

- Backend: AWS Lambda, API Gateway HTTP API, DynamoDB e CloudWatch Logs.
- Frontend: Vercel.

Siga o guia da AWS do início ao fim: ele inclui pacote ARM64, IAM com privilégio
mínimo, variáveis, rotas, CORS, testes e diagnóstico.

## CI do backend (v1)

O repositório possui uma pipeline inicial de backend em GitHub Actions
([`.github/workflows/backend-ci.yml`](.github/workflows/backend-ci.yml)).

Ela roda em `push` e `pull_request` da branch `main`, com escopo apenas em
mudanças de `backend/**`, e executa:

- validação de formatação (`go fmt ./...`);
- análise estática (`go vet ./...`);
- testes (`make test`);
- verificação de dependências (`go mod tidy` sem diff pendente);
- build e empacotamento da Lambda (`make package`), publicando
  `backend/bin/function.zip` como artifact do workflow.

Nesta primeira versão, o deploy na AWS continua manual.

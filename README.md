# Todo API com AWS Lambda (Go)

API REST simples de lista de tarefas em Go, executando como uma unica AWS Lambda Function e recebendo requisicoes via API Gateway HTTP API (Payload Format Version 2.0).

## Arquitetura resumida

- API Gateway HTTP API encaminha todas as rotas para uma unica Lambda.
- A Lambda identifica metodo e caminho e processa as rotas manualmente.
- Os dados sao armazenados em uma tabela DynamoDB.

## Rotas

- `GET /todos`: lista todas as tarefas.
- `POST /todos`: cria uma tarefa.
- `PATCH /todos/{id}`: atualiza campos enviados (`title` e/ou `completed`).
- `DELETE /todos/{id}`: remove tarefa permanentemente.

Todas as respostas com body usam JSON e incluem os headers:

- `Content-Type: application/json`
- `Access-Control-Allow-Origin: *`

## Variavel de ambiente

- `TABLE_NAME`: nome da tabela DynamoDB.

## Estrutura da tabela DynamoDB

- Tabela com partition key:
- `id` (String)

Modelo salvo:

```json
{
  "id": "uuid",
  "title": "Estudar AWS",
  "completed": false,
  "createdAt": "2026-07-26T20:00:00Z",
  "updatedAt": "2026-07-26T20:00:00Z"
}
```

## Build e empacotamento

```bash
make build
make package
make clean
```

- `make build` gera o executavel `bin/bootstrap` para Linux arm64.
- `make package` gera `function.zip` com `bootstrap` na raiz do zip.

## Quickstart local

1. Copie o arquivo de ambiente:

```bash
cp .env.example .env
```

2. (Opcional) Edite o `.env`:
   - `TABLE_NAME`
   - `AWS_REGION`
  - `DYNAMODB_ENDPOINT`

3. Suba o servidor e o DynamoDB local:

```bash
make start
```

4. Confira logs/status:

```bash
make logs
make ps
```

Comandos de ciclo de vida:

- `make restart`: reinicia o ambiente local com rebuild da imagem.
- `make stop`: derruba o ambiente local.

Observacoes:

- O `make start` sobe apenas o servidor HTTP e o LocalStack com DynamoDB.
- O container do runtime Lambda nao e necessario no desenvolvimento local.
- A tabela DynamoDB e criada automaticamente no boot do LocalStack.
- Nao e necessario criar nada na AWS para rodar local.
- As imagens Docker funcionam em Raspberry Pi, Windows, macOS e Linux (build por arquitetura automaticamente).

## Teste rapido

URL:

- `POST http://localhost:8081/todos`

Header:

- `Content-Type: application/json`

Body:

```json
{
  "title": "Estudar Go"
}
```

As demais rotas estao disponiveis diretamente em `http://localhost:8081`.

## Configuracao da Lambda

- Runtime: `provided.al2023`
- Architecture: `arm64`
- Handler: `bootstrap`
- Variavel de ambiente: `TABLE_NAME`

No API Gateway HTTP API, use integracao com Lambda proxy e Payload Format Version 2.0.

## Permissoes IAM minimas da Lambda

- `dynamodb:Scan`
- `dynamodb:PutItem`
- `dynamodb:GetItem`
- `dynamodb:UpdateItem`
- `dynamodb:DeleteItem`

A role da Lambda tambem precisa de permissao para gravar logs no CloudWatch.

# Rodar localmente e testar com Postman

O projeto usa um servidor HTTP local que reaproveita os mesmos handlers da
AWS Lambda. O DynamoDB e executado pelo LocalStack.

## 1. Pre-requisitos

- Docker

## 2. Configurar variaveis de ambiente

Na raiz do projeto:

```bash
cp .env.example .env
```

Depois, edite o `.env` se necessario:

- `TABLE_NAME`
- `AWS_REGION`
- `DYNAMODB_ENDPOINT`

## 3. Subir o ambiente

```bash
make start
```

A tabela DynamoDB e criada automaticamente no LocalStack.

## 4. Testar no Postman

### Criar tarefa

- Metodo: `POST`
- URL: `http://localhost:8081/todos`
- Header: `Content-Type: application/json`
- Body (raw JSON):

```json
{
  "title": "Estudar Go"
}
```

### Listar tarefas

- Metodo: `GET`
- URL: `http://localhost:8081/todos`

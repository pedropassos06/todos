# Backend

API REST em Go com DynamoDB. O mesmo pacote de handlers é usado pelo servidor
HTTP local e pela função AWS Lambda.

## Arquitetura de produção

- AWS Lambda com runtime `provided.al2023` e arquitetura ARM64;
- API Gateway HTTP API com payload format `2.0`;
- DynamoDB com `id` String como partition key;
- CloudWatch Logs.

O ponto de entrada local é `cmd/server`; o da Lambda é `cmd/lambda`.

## Executar localmente

Na raiz do projeto:

```bash
cp .env.example .env
cd backend
make start
```

A API fica disponível em `http://localhost:8081`.

O arquivo só precisa ser copiado na primeira execução. `make start` inicia a
API e o DynamoDB no LocalStack.

## Comandos

Execute os comandos abaixo dentro de `backend/`:

```bash
make start    # Start the API and DynamoDB
make stop     # Stop the local services
make logs     # Follow API logs
make test     # Run Go tests
make build    # Build the AWS Lambda executable
make package  # Generate function.zip for AWS Lambda
make clean    # Remove backend build artifacts
```

A Lambda usa somente `TABLE_NAME` e `ALLOWED_ORIGIN` como variáveis próprias
da aplicação em produção. `AWS_REGION`, credenciais fictícias e
`DYNAMODB_ENDPOINT` no `.env` da raiz são para o LocalStack.

## Rotas

| Método | Caminho | Sucesso |
| --- | --- | --- |
| `GET` | `/todos` | `200` |
| `POST` | `/todos` | `201` |
| `PATCH` | `/todos/{id}` | `200` |
| `DELETE` | `/todos/{id}` | `204` |
| `OPTIONS` | qualquer rota configurada no API Gateway | `204` |

O contrato completo e os exemplos estão em
[`../docs/rodar-local-postman.md`](../docs/rodar-local-postman.md).

O deploy manual na AWS está em
[`../docs/deploy-manual-aws-console.md`](../docs/deploy-manual-aws-console.md).

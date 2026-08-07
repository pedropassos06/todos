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
cd backend
make start
```

A API fica disponível em `http://localhost:8081`.

Na primeira execução, `make start` cria `backend/.env` automaticamente a
partir de `backend/.env.example` e inicia a API com o DynamoDB no LocalStack.

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
`DYNAMODB_ENDPOINT` no `backend/.env` são para o LocalStack.

No `.env` local, você pode usar `ALLOWED_ORIGIN=*` para evitar bloqueios de
CORS em diferentes hosts de desenvolvimento. Em produção, configure a origem
exata do frontend.

## Rotas

| Método | Caminho | Sucesso |
| --- | --- | --- |
| `GET` | `/todos` | `200` |
| `POST` | `/todos` | `201` |
| `PATCH` | `/todos/{id}` | `200` |
| `DELETE` | `/todos/{id}` | `204` |
| `OPTIONS` | qualquer rota configurada no API Gateway | `204` |

O contrato completo e os exemplos estão em
[`./docs/rodar-local-postman.md`](./docs/rodar-local-postman.md).

O deploy manual na AWS está em
[`./docs/deploy-manual-aws-console.md`](./docs/deploy-manual-aws-console.md).

## CI/CD

A pipeline inicial do backend está em
[`../.github/workflows/backend-ci.yml`](../.github/workflows/backend-ci.yml)
e roda em `push` e `pull_request` da branch `main`, apenas quando há mudanças
no backend.

Checks executados:

- `go fmt ./...` (falha se houver arquivo não formatado);
- `go vet ./...`;
- `make test`;
- `go mod tidy` com verificação de diff limpo;
- `make package` para gerar `function.zip` como artifact.

O deploy automatizado do backend está em
[`./docs/deploy-github-actions.md`](./docs/deploy-github-actions.md).

O guia manual no Console AWS continua disponível como fallback em
[`./docs/deploy-manual-aws-console.md`](./docs/deploy-manual-aws-console.md).

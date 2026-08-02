# Rodar localmente e testar com Postman

O ambiente local executa o servidor Go e um DynamoDB no LocalStack. O servidor
converte cada requisição HTTP para o mesmo evento API Gateway v2 usado pela
Lambda, portanto exercita os mesmos handlers do deploy AWS.

## 1. Pré-requisitos

- Docker com Docker Compose;
- portas `4566` e `8081` livres.

Go só é necessário para executar `make test` ou compilar fora do Docker.

## 2. Configurar o ambiente

Na raiz do projeto (`todos/`):

```bash
cp backend/.env.example backend/.env
```

O arquivo de exemplo já contém os valores locais:

```dotenv
TABLE_NAME=todos-table
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
DYNAMODB_ENDPOINT=http://localstack:4566
ALLOWED_ORIGIN=http://localhost:5173
```

As credenciais `test` são fictícias e servem apenas para o LocalStack. Não use
esses valores no deploy AWS.

Essas variáveis do `backend/.env.example` atendem apenas ao backend no ambiente local.

## 3. Subir e conferir os serviços

Dentro de `backend/`:

```bash
cd backend
make start
```

A tabela DynamoDB e criada automaticamente no LocalStack.

## 4. Testar no Postman

### Criar tarefa

- Método: `POST`
- URL: `http://localhost:8081/todos`
- Header: `Content-Type: application/json`
- Body: **raw > JSON**

```json
{
  "title": "Estudar Go"
}
```

Isso deve retornar com sucesso e código 201.

### Listar tarefas

- Método: `GET`
- URL: `http://localhost:8081/todos`

Resposta esperada: `200 OK`.

Quando a tabela está vazia, a resposta é `{"todos":[]}`.


## Parar ou recriar o ambiente

```bash
cd backend
make stop
```

O `docker compose down` remove os containers e a rede. Como o projeto não
configura volume persistente para o LocalStack, os dados locais são descartados
quando o serviço é removido.

Para subir tudo novamente:

```bash
cd backend
make start
```

## Problemas comuns

- Porta ocupada: libere `4566` ou `8081` antes de executar `make start`.
- API retorna `500`: consulte `make logs` e confira o `.env`.
- Tabela não existe: execute `make stop` e `make start` para rodar novamente o
  script de inicialização do LocalStack.
- Frontend bloqueado por CORS: mantenha
  `ALLOWED_ORIGIN=http://localhost:5173` no ambiente local.

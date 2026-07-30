# Deploy do backend no Console da AWS

Este guia publica o backend no estado atual do projeto:

```text
cliente -> API Gateway HTTP API -> Lambda Go -> DynamoDB
                                      |
                                      +-> CloudWatch Logs
```

A API inteira roda em uma única Lambda. O procedimento abaixo usa o Console da
AWS e o CORS implementado pelo próprio backend.

## 1. Antes de começar

Você precisa de:

- uma conta AWS com permissão para administrar DynamoDB, IAM, Lambda e API
  Gateway;
- Go 1.24, `make` e `zip` instalados na máquina;
- o repositório clonado.

Escolha uma região e crie todos os recursos nela. Os exemplos usam
`us-east-1`. Anote estes valores:

```text
REGIAO=us-east-1
NOME_TABELA=todos-table
NOME_ROLE=todos-lambda-role
NOME_LAMBDA=todos-api
ORIGEM_FRONTEND=https://seu-projeto.vercel.app
```

`ORIGEM_FRONTEND` é somente a origem, sem caminho e sem `/` no final. Se o
frontend ainda não foi publicado, use temporariamente `*` e restrinja o valor
depois do deploy do frontend.

> A AWS cobra pelos recursos usados. Para este projeto, DynamoDB sob demanda,
> Lambda e HTTP API são as opções mais simples, mas continuam sujeitos à
> cobrança da conta.

## 2. Gerar e conferir o pacote da Lambda

Na raiz do projeto (`todos/`):

```bash
make clean
make test
make package
unzip -l backend/function.zip
```

O último comando deve mostrar um arquivo chamado `bootstrap` diretamente na
raiz do ZIP. O `Makefile` compila para Linux ARM64, portanto a Lambda também
precisa usar a arquitetura `arm64`.

Artefato que será enviado:

```text
backend/function.zip
```

## 3. Criar a tabela DynamoDB

No Console da AWS:

1. Abra **DynamoDB** e confirme a região escolhida.
2. Entre em **Tables > Create table**.
3. Use:
   - **Table name:** `todos-table` ou o valor de `NOME_TABELA`;
   - **Partition key:** `id`;
   - **Partition key type:** `String`;
   - **Sort key:** deixe desmarcada;
   - **Table settings:** `Default settings`.
4. Crie a tabela e aguarde o status **Active**.

As configurações padrão usam capacidade sob demanda. Não é preciso criar
índices nem inserir itens manualmente.

Na página da tabela, copie o **Amazon Resource Name (ARN)**. Ele terá este
formato:

```text
arn:aws:dynamodb:REGIAO:ID_DA_CONTA:table/NOME_TABELA
```

## 4. Criar a role de execução da Lambda

### 4.1 Criar a role

1. Abra **IAM > Roles > Create role**.
2. Em **Trusted entity type**, escolha `AWS service`.
3. Em **Use case**, escolha `Lambda`.
4. Anexe a policy gerenciada `AWSLambdaBasicExecutionRole`.
5. Nomeie a role como `todos-lambda-role` ou o valor de `NOME_ROLE` e conclua.

Essa policy permite enviar logs ao CloudWatch. Ela não dá acesso aos itens do
DynamoDB.

### 4.2 Adicionar acesso somente à tabela do projeto

Abra a role criada e selecione **Add permissions > Create inline policy**.
Na aba **JSON**, use o documento abaixo, substituindo `REGIAO`,
`ID_DA_CONTA` e `NOME_TABELA` pelo ARN copiado no passo anterior:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "TodosTableAccess",
      "Effect": "Allow",
      "Action": [
        "dynamodb:Scan",
        "dynamodb:PutItem",
        "dynamodb:GetItem",
        "dynamodb:UpdateItem",
        "dynamodb:DeleteItem"
      ],
      "Resource": "arn:aws:dynamodb:REGIAO:ID_DA_CONTA:table/NOME_TABELA"
    }
  ]
}
```

Avance, nomeie a policy como `todos-table-access` e salve.

Não use `AWSLambdaDynamoDBExecutionRole` como substituta: essa policy é voltada
a consumir DynamoDB Streams e não concede as operações CRUD usadas por esta
API.

## 5. Criar a Lambda

1. Abra **Lambda > Functions > Create function**.
2. Escolha **Author from scratch**.
3. Configure:
   - **Function name:** `todos-api` ou o valor de `NOME_LAMBDA`;
   - **Runtime:** `Amazon Linux 2023` / `provided.al2023`;
   - **Architecture:** `arm64`;
   - **Execution role:** `Use an existing role`;
   - **Existing role:** a role criada no passo 4.
4. Crie a função.

Na aba **Code**:

1. Escolha **Upload from > .zip file**.
2. Envie `backend/function.zip`.
3. Clique em **Save** se o Console solicitar.

Em **Code > Runtime settings > Edit**, confirme:

```text
Handler: bootstrap
```

O runtime exige que o executável `bootstrap` esteja na raiz do ZIP. O nome do
handler não é usado para localizar uma função Go específica, mas mantê-lo como
`bootstrap` evita ambiguidade.

## 6. Configurar a Lambda

Em **Configuration > Environment variables > Edit**, adicione:

| Chave | Valor |
| --- | --- |
| `TABLE_NAME` | nome exato da tabela, por exemplo `todos-table` |
| `ALLOWED_ORIGIN` | origem exata do frontend, por exemplo `https://seu-projeto.vercel.app` |

Não configure `DYNAMODB_ENDPOINT`, `AWS_ACCESS_KEY_ID` nem
`AWS_SECRET_ACCESS_KEY` na Lambda. Esses valores do `.env` existem somente
para o ambiente local com LocalStack. A região e as credenciais em produção
são fornecidas pelo ambiente e pela role de execução da Lambda.

Em **Configuration > General configuration > Edit**, uma configuração inicial
adequada é:

```text
Memory: 128 MB
Timeout: 10 seconds
```

Salve e aguarde a atualização terminar.

## 7. Criar a HTTP API

1. Abra **API Gateway > APIs > Create API**.
2. Em **HTTP API**, escolha **Build**. Não escolha REST API.
3. Em **Integrations**, selecione `Lambda`.
4. Escolha a região e a Lambda `todos-api`.
5. Dê um nome à API, por exemplo `todos-http-api`.
6. Adicione estas cinco rotas, todas apontando para a mesma integração:

| Método | Recurso |
| --- | --- |
| `GET` | `/todos` |
| `POST` | `/todos` |
| `PATCH` | `/todos/{id}` |
| `DELETE` | `/todos/{id}` |
| `OPTIONS` | `/{proxy+}` |

7. Use o stage `$default` com **Auto-deploy** habilitado.
8. Revise e crie a API.

Abra **Integrations**, selecione a integração da Lambda e confirme:

```text
Payload format version: 2.0
```

O backend usa o evento `APIGatewayV2HTTPRequest`, portanto o payload `2.0` é
obrigatório.

### CORS: configuração usada por este projeto

Em **CORS** no API Gateway, deixe o CORS nativo **desativado**. O backend já
retorna os cabeçalhos:

```text
Access-Control-Allow-Origin
Access-Control-Allow-Headers
Access-Control-Allow-Methods
```

A rota `OPTIONS /{proxy+}` encaminha o preflight à Lambda. Se o CORS nativo do
HTTP API for habilitado, o API Gateway passa a ignorar os cabeçalhos CORS
retornados pela Lambda e `ALLOWED_ORIGIN` deixa de controlar a origem efetiva.
Use uma estratégia ou a outra, não as duas.

Ao criar a integração pelo Console, o API Gateway normalmente adiciona à
Lambda a permissão para invocá-la. Em **Lambda > Configuration >
Permissions > Resource-based policy statements**, confirme que há uma
permissão com principal `apigateway.amazonaws.com`.

## 8. Testar o backend publicado

Na página da API, copie o **Invoke URL** do stage `$default`. Ele se parece
com:

```text
https://abc123.execute-api.us-east-1.amazonaws.com
```

Defina a URL no terminal:

```bash
API_URL="https://abc123.execute-api.us-east-1.amazonaws.com"
```

Não acrescente o nome do stage quando estiver usando `$default`.

### Criar

```bash
curl --fail-with-body -i -X POST "$API_URL/todos" \
  -H "Content-Type: application/json" \
  -d '{"title":"Validar deploy na AWS"}'
```

Espere `HTTP 201`. Copie o campo `id` retornado.

### Listar

```bash
curl --fail-with-body -i "$API_URL/todos"
```

Espere `HTTP 200` e um JSON no formato `{"todos":[...]}`.

### Atualizar

```bash
TODO_ID="ID_RETORNADO_NO_POST"

curl --fail-with-body -i -X PATCH "$API_URL/todos/$TODO_ID" \
  -H "Content-Type: application/json" \
  -d '{"title":"Deploy validado","completed":true}'
```

Espere `HTTP 200`.

### Excluir

```bash
curl --fail-with-body -i -X DELETE "$API_URL/todos/$TODO_ID"
```

Espere `HTTP 204`.

### Validar o preflight CORS

Use exatamente a origem configurada em `ALLOWED_ORIGIN`:

```bash
FRONTEND_ORIGIN="https://seu-projeto.vercel.app"

curl --fail-with-body -i -X OPTIONS "$API_URL/todos" \
  -H "Origin: $FRONTEND_ORIGIN" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type"
```

Espere `HTTP 204` e, entre os cabeçalhos:

```text
Access-Control-Allow-Origin: https://seu-projeto.vercel.app
Access-Control-Allow-Headers: Content-Type
Access-Control-Allow-Methods: GET,POST,PATCH,DELETE,OPTIONS
```

## 9. Ligar o frontend ao backend

No deploy Vercel descrito em `frontend/README.md`, configure:

```text
VITE_API_BASE_URL=https://abc123.execute-api.us-east-1.amazonaws.com
```

Faça um novo deploy do frontend após alterar a variável. Em seguida, volte à
Lambda e confira que:

```text
ALLOWED_ORIGIN=https://URL-EXATA-DO-FRONTEND
```

Não coloque `/todos` em `VITE_API_BASE_URL`: o frontend acrescenta esse
caminho.

## 10. Diagnóstico

### `Internal Server Error` ou resposta `500`

Abra **CloudWatch > Log groups > `/aws/lambda/todos-api`** e consulte o stream
mais recente.

- `TABLE_NAME environment variable is required`: crie a variável na Lambda.
- `ResourceNotFoundException`: confira o nome e a região da tabela.
- `AccessDeniedException`: confira a inline policy e o ARN da tabela.
- timeout: aumente temporariamente o timeout e confirme que a Lambda não foi
  colocada em uma VPC sem acesso aos endpoints da AWS.

### `Not Found` ou resposta `404`

Confira método e caminho nas rotas. O código reconhece exatamente:

```text
GET /todos
POST /todos
PATCH /todos/{id}
DELETE /todos/{id}
```

Confirme também que o stage `$default` está com auto-deploy ativo.

### Erro de CORS no navegador

1. Execute o teste de preflight do passo 8.
2. Confira se `ALLOWED_ORIGIN` tem somente esquema e domínio, sem `/` final.
3. Confirme que existe `OPTIONS /{proxy+}` integrado à Lambda.
4. Confirme que o CORS nativo do HTTP API está desativado.
5. Depois de alterar variáveis da Lambda, recarregue o frontend sem cache.

### A Lambda não inicia (`Runtime.InvalidEntrypoint`)

Confirme os três itens em conjunto:

- ZIP gerado por `make package`;
- arquitetura da Lambda `arm64`;
- arquivo `bootstrap` na raiz do ZIP.

## 11. Atualizar uma Lambda existente

Depois de alterar o backend:

```bash
make test
make clean
make package
```

Envie novamente `backend/function.zip` em **Lambda > Code > Upload from >
.zip file**. Não é necessário recriar a tabela, a role ou a API.

## 12. Checklist final

- [ ] todos os recursos estão na mesma região;
- [ ] tabela ativa, com partition key `id` do tipo String e sem sort key;
- [ ] role com `AWSLambdaBasicExecutionRole` e a inline policy da tabela;
- [ ] Lambda `provided.al2023`, `arm64`, handler `bootstrap`;
- [ ] `TABLE_NAME` e `ALLOWED_ORIGIN` configuradas;
- [ ] HTTP API, payload `2.0`, stage `$default` e auto-deploy;
- [ ] cinco rotas criadas, incluindo `OPTIONS /{proxy+}`;
- [ ] CORS nativo do API Gateway desativado;
- [ ] testes de criar, listar, atualizar, excluir e preflight concluídos;
- [ ] frontend usando o Invoke URL sem `/todos` no final.

## Referências oficiais da AWS

- [Empacotar funções Lambda em Go com arquivo ZIP](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)
- [Integração proxy entre HTTP API e Lambda](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-integrations-lambda.html)
- [Rotas de HTTP APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-routes.html)
- [CORS em HTTP APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-cors.html)
- [Role de execução da Lambda](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html)
- [Logs da Lambda no CloudWatch](https://docs.aws.amazon.com/lambda/latest/dg/monitoring-cloudwatchlogs.html)

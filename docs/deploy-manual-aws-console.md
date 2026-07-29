# Deploy Manual no Console da AWS

Este guia mostra o passo a passo para publicar esta API de tarefas em uma unica AWS Lambda Function, usando API Gateway HTTP API e DynamoDB.

## 1. Preparar o pacote localmente

Na raiz do projeto, execute:

```bash
make clean
make build
make package
```

No final, o arquivo `function.zip` deve existir na raiz.

## 2. Criar a tabela DynamoDB

No Console AWS:

1. Acesse DynamoDB.
2. Clique em Create table.
3. Configure:
   - Table name: o nome que voce quiser (exemplo: `todos-table`)
   - Partition key: `id`
   - Type: `String`
4. Clique em Create table.

## 3. Criar a Role IAM da Lambda

No Console AWS:

1. Acesse IAM > Roles > Create role.
2. Escolha AWS service > Lambda.
3. Adicione permissoes minimas para DynamoDB:
   - `dynamodb:Scan`
   - `dynamodb:PutItem`
   - `dynamodb:GetItem`
   - `dynamodb:UpdateItem`
   - `dynamodb:DeleteItem`
4. Garanta permissao de logs no CloudWatch (exemplo: policy AWSLambdaBasicExecutionRole).
5. Finalize a criacao da role.

## 4. Criar a Lambda Function

No Console AWS:

1. Acesse Lambda > Create function.
2. Escolha Author from scratch.
3. Configure:
   - Function name: exemplo `todos-api`
   - Runtime: `provided.al2023`
   - Architecture: `arm64`
   - Execution role: selecione a role criada no passo anterior
4. Clique em Create function.

## 5. Enviar o function.zip

Na tela da Lambda:

1. Abra a aba Code.
2. Clique em Upload from > .zip file.
3. Envie `function.zip`.
4. Em Runtime settings, confirme:
   - Handler: `bootstrap`

## 6. Configurar variavel de ambiente

Na Lambda:

1. Abra Configuration > Environment variables.
2. Clique em Edit > Add environment variable.
3. Configure:
   - Key: `TABLE_NAME`
   - Value: nome exato da tabela criada no DynamoDB
4. Salve.

## 7. Criar API Gateway HTTP API

No Console AWS:

1. Acesse API Gateway > Create API > HTTP API.
2. Em Integrations, adicione integracao com a Lambda criada.
3. Crie as rotas:
   - `GET /todos`
   - `POST /todos`
   - `PATCH /todos/{id}`
   - `DELETE /todos/{id}`
4. Associe todas as rotas a mesma Lambda.
5. Em Payload format version, use `2.0`.
6. Crie ou use o stage `$default`.

## 8. Testar a API

Copie a URL base da HTTP API (exemplo: `https://abc123.execute-api.us-east-1.amazonaws.com`).

### Criar tarefa

```bash
curl -X POST "https://SUA-URL/todos" \
  -H "Content-Type: application/json" \
  -d '{"title":"Estudar AWS"}'
```

### Listar tarefas

```bash
curl "https://SUA-URL/todos"
```

### Atualizar tarefa

```bash
curl -X PATCH "https://SUA-URL/todos/SEU_ID" \
  -H "Content-Type: application/json" \
  -d '{"title":"Estudar AWS Lambda","completed":true}'
```

### Excluir tarefa

```bash
curl -X DELETE "https://SUA-URL/todos/SEU_ID"
```

## 9. Checklist rapido

- `function.zip` foi gerado localmente
- tabela DynamoDB criada com `id` como partition key String
- Lambda criada com `provided.al2023` e `arm64`
- handler configurado como `bootstrap`
- variavel `TABLE_NAME` configurada
- role IAM com permissoes DynamoDB e CloudWatch Logs
- API Gateway HTTP API com rotas e payload `2.0`

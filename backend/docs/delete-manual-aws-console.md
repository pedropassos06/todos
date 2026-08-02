# Remover o backend criado manualmente no Console da AWS

Este guia remove os recursos criados pelo procedimento
`deploy-manual-aws-console.md`:

```text
API Gateway HTTP API -> Lambda Go -> DynamoDB
                            |
                            +-> CloudWatch Logs
                            |
                            +-> IAM execution role
```

Siga a ordem indicada para remover primeiro os recursos que dependem dos
demais. As exclusões são permanentes e interrompem imediatamente o backend.

## 1. Antes de começar

Confirme que você está na mesma conta e região usadas no deploy. Os exemplos
do guia de criação usam:

```text
REGIAO=us-east-1
NOME_API=todos-http-api
NOME_LAMBDA=todos-api
NOME_ROLE=todos-lambda-role
NOME_TABELA=todos-table
LOG_GROUP=/aws/lambda/todos-api
```

Se você usou nomes diferentes, substitua-os em todos os passos. Antes de
excluir qualquer recurso, confira seu nome, ARN, região e, quando exibido, o ID
da conta. Não remova um recurso apenas porque o nome é parecido.

> A exclusão da tabela apaga todos os itens. Se os dados precisarem ser
> preservados, crie um backup antes de continuar.

## 2. Excluir a HTTP API

Remova primeiro a entrada pública do backend:

1. Abra **API Gateway > APIs** e confirme a região.
2. Localize `todos-http-api` ou o nome usado no deploy.
3. Abra a API e confira se as rotas são as esperadas:
   - `GET /todos`;
   - `POST /todos`;
   - `PATCH /todos/{id}`;
   - `DELETE /todos/{id}`;
   - `OPTIONS /{proxy+}`.
4. Volte à lista de APIs, selecione somente essa API e escolha **Delete**.
5. Confirme a exclusão quando solicitado.

Depois disso, o antigo Invoke URL deve deixar de responder. Não é necessário
excluir separadamente as rotas, a integração ou o stage `$default`: eles fazem
parte da API.

## 3. Excluir a função Lambda

1. Abra **Lambda > Functions** e confirme a região.
2. Selecione `todos-api`.
3. Escolha **Actions > Delete function**.
4. Digite a confirmação solicitada pelo Console e exclua a função.

A exclusão da função remove seu código, configurações, variáveis de ambiente e
resource policy. Ela não exclui automaticamente o log group do CloudWatch nem
a role de execução.

## 4. Excluir os logs do CloudWatch

Este passo remove o histórico de logs e evita que ele permaneça armazenado:

1. Abra **CloudWatch > Logs > Log Management** (sempre confira a região).
2. Localize `/aws/lambda/todos-api`.
3. Selecione somente esse log group.
4. Escolha **Actions > Delete log group(s)** e confirme.

Se o log group não existir, nenhuma ação é necessária. Não exclua outros log
groups com o prefixo `/aws/lambda/`.

## 5. Excluir a role de execução do IAM

IAM é um serviço global; mesmo assim, confirme que a role é a usada pela
função removida.

1. Abra **IAM > Roles**.
2. Selecione `todos-lambda-role`.
3. Verifique que esse é exatamente o role que você criou incialmente.
4. Escolha **Delete**, digite o nome da role e confirme.

## 6. Excluir a tabela DynamoDB

Faça esta etapa por último, pois ela apaga permanentemente os dados:

1. Abra **DynamoDB > Tables** e confirme a região.
2. Selecione `todos-table`.
3. Escolha **Delete**.
4. Selecione para deletar todos os alarmes CloudWatch.
5. Não precisa selecionar a opção de **backup**, a não ser que você realmente precise.
6. Digite a confirmação solicitada e exclua a tabela.
7. Aguarde até que ela desapareça da lista de tabelas.

Se a proteção contra exclusão estiver habilitada, abra as configurações da
tabela, desabilite **Deletion protection**, salve a alteração e tente
novamente.

## 7. Opcional: remover a referência no frontend

Este passo não exclui recursos da AWS. Ele apenas evita que o frontend continue
chamando uma API que não existe mais.

No projeto da Vercel usado no deploy:

1. Abra **Settings > Environment Variables**.
2. Remova `VITE_API_BASE_URL` ou substitua seu valor pelo endereço do novo
   backend.
3. Faça um novo deploy do frontend para aplicar a alteração.

Não exclua o projeto da Vercel a menos que o frontend também deva ser
desativado.

## 8. Verificação final

Confira no Console, usando a conta e região corretas:

- **API Gateway > APIs:** `todos-http-api` não aparece;
- **Lambda > Functions:** `todos-api` não aparece;
- **CloudWatch > Log groups:** `/aws/lambda/todos-api` não aparece;
- **IAM > Roles:** `todos-lambda-role` não aparece, exceto se foi preservada
  por ser compartilhada;
- **DynamoDB > Tables:** `todos-table` não aparece.

Se você criou um backup no passo 2, ele continuará listado em
**DynamoDB > Backups**. Mantenha-o apenas pelo tempo necessário e exclua-o
manualmente quando não precisar mais dos dados.

O arquivo local `backend/function.zip` não é um recurso da AWS e não é afetado
por este procedimento. Para removê-lo e limpar outros artefatos locais de
build, execute os comandos abaixo em `backend/`:

```bash
cd backend
make clean
```

# Deploy automatizado do backend com GitHub Actions

Este guia publica e atualiza a infraestrutura do backend na AWS usando o
workflow de CI/CD do repositório.

Fluxo automatizado:

```text
push na main -> GitHub Actions -> build backend/bin/function.zip -> update Lambda -> CloudFormation deploy
```

A stack usa o template em `backend/infra/cloudformation.yaml`.

## 1. O que o workflow faz

No push para `main` (com mudanças em `backend/**`):

1. valida o backend (`go fmt`, `go vet`, testes, `go mod tidy`);
2. gera `backend/bin/function.zip`;
3. atualiza o código da Lambda com `aws lambda update-function-code` usando o zip local;
4. atualiza configuração da Lambda (memória, timeout e variáveis de ambiente);
5. executa `aws cloudformation deploy` com os parâmetros da stack;
5. publica um resumo com os outputs (`ApiEndpoint`, `LambdaFunctionName`, `DynamoDBTableName`).

Não há uso de bucket S3 nesse fluxo.

## 2. Pré-requisitos

- Conta AWS com permissão para criar/atualizar Lambda, API Gateway HTTP API,
  DynamoDB, IAM e CloudWatch Logs via CloudFormation.
- Repositório no GitHub com Actions habilitado.
- Branch principal chamada `main`.
- A função Lambda definida em `FUNCTION_NAME` já deve existir antes do primeiro
  deploy automático (o workflow só atualiza código/configuração da função).

## 3. Criar credencial AWS para o GitHub Actions

Este projeto está configurado para usar secrets com chave estática:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

### 4.1 Criar usuário IAM para CI

No IAM:

1. Crie um usuário (exemplo: `github-actions-todos-backend`).
2. Habilite **Access key** para uso programático.
3. Guarde `Access key ID` e `Secret access key`.

### 4.2 Permissões IAM mínimas (policy exemplo)

Anexe uma policy ao usuário com escopo na sua região/conta e, quando possível,
na stack específica.

Exemplo inicial funcional (ajuste `REGION`, `ACCOUNT_ID` e `STACK_NAME`):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "CloudFormationDeploy",
      "Effect": "Allow",
      "Action": [
        "cloudformation:CreateStack",
        "cloudformation:UpdateStack",
        "cloudformation:DescribeStacks",
        "cloudformation:DescribeStackEvents",
        "cloudformation:DescribeStackResources",
        "cloudformation:GetTemplateSummary",
        "cloudformation:ValidateTemplate",
        "cloudformation:CreateChangeSet",
        "cloudformation:DescribeChangeSet",
        "cloudformation:ExecuteChangeSet",
        "cloudformation:DeleteChangeSet",
        "cloudformation:ListStackResources"
      ],
      "Resource": [
        "arn:aws:cloudformation:REGION:ACCOUNT_ID:stack/STACK_NAME/*"
      ]
    },
    {
      "Sid": "AllowCreateStackWhenMissing",
      "Effect": "Allow",
      "Action": [
        "cloudformation:CreateStack",
        "cloudformation:DescribeStacks",
        "cloudformation:GetTemplateSummary",
        "cloudformation:ValidateTemplate"
      ],
      "Resource": "*"
    },
    {
      "Sid": "PassRolesCreatedByStack",
      "Effect": "Allow",
      "Action": "iam:PassRole",
      "Resource": "arn:aws:iam::ACCOUNT_ID:role/STACK_NAME-*"
    },
    {
      "Sid": "ManageStackResources",
      "Effect": "Allow",
      "Action": [
        "iam:CreateRole",
        "iam:DeleteRole",
        "iam:GetRole",
        "iam:PutRolePolicy",
        "iam:DeleteRolePolicy",
        "iam:AttachRolePolicy",
        "iam:DetachRolePolicy",
        "iam:TagRole",
        "iam:UntagRole",
        "iam:CreatePolicy",
        "iam:DeletePolicy",
        "iam:GetPolicy",
        "iam:GetPolicyVersion",
        "iam:CreatePolicyVersion",
        "iam:DeletePolicyVersion",
        "lambda:CreateFunction",
        "lambda:DeleteFunction",
        "lambda:GetFunction",
        "lambda:GetFunctionConfiguration",
        "lambda:UpdateFunctionCode",
        "lambda:UpdateFunctionConfiguration",
        "lambda:AddPermission",
        "lambda:RemovePermission",
        "apigateway:POST",
        "apigateway:GET",
        "apigateway:PUT",
        "apigateway:PATCH",
        "apigateway:DELETE",
        "dynamodb:CreateTable",
        "dynamodb:DeleteTable",
        "dynamodb:DescribeTable",
        "dynamodb:UpdateTable",
        "logs:CreateLogGroup",
        "logs:DeleteLogGroup",
        "logs:PutRetentionPolicy",
        "logs:DeleteRetentionPolicy",
        "logs:DescribeLogGroups",
        "tag:GetResources"
      ],
      "Resource": "*"
    }
  ]
}
```

Observação: o bloco `ManageStackResources` está amplo para reduzir chance de
falha no primeiro deploy. Depois do primeiro deploy estável, você pode reduzir
escopo com base no CloudTrail.

## 4. Configurar Secrets e Variables no GitHub

No repositório, abra Settings > Secrets and variables > Actions.

### 5.1 Repository secrets

Crie:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

### 5.2 Repository variables

Crie:

- `AWS_REGION=us-east-1`
- `CFN_STACK_NAME=todos-backend-prod`
- `TABLE_NAME=todos-table`
- `FUNCTION_NAME=todos-api`
- `ALLOWED_ORIGIN=https://seu-frontend.vercel.app`
- `LAMBDA_MEMORY=128`
- `LAMBDA_TIMEOUT=10`

Notas:

- Em produção, não use `ALLOWED_ORIGIN=*`; use origem exata do frontend.
- Não inclua `/todos` em `ALLOWED_ORIGIN`.

## 5. Primeiro deploy

1. Faça merge de uma alteração em `main` que toque `backend/**`.
2. Abra a aba Actions e acompanhe o workflow **Backend CI**.
3. No job `Deploy Backend (AWS)`, confirme:
  - atualização do código/configuração da Lambda bem-sucedida;
   - `cloudformation deploy` concluído.
4. Leia o resumo do job para obter `ApiEndpoint`.

## 6. Teste pós-deploy

Teste rápido com `curl`:

```bash
API_URL="https://<api-id>.execute-api.us-east-1.amazonaws.com"
curl -i "${API_URL}/todos"
```

Você deve receber `200 OK`.

Teste preflight CORS:

```bash
curl -i -X OPTIONS "${API_URL}/todos" \
  -H "Origin: https://seu-frontend.vercel.app" \
  -H "Access-Control-Request-Method: GET"
```

Confirme presença de:

- `Access-Control-Allow-Origin` com a origem esperada;
- `Access-Control-Allow-Methods`;
- `Access-Control-Allow-Headers`.

## 7. Atualizações de configuração

Para alterar memória, timeout, nome de tabela, nome da função ou CORS:

1. atualize a repository variable correspondente;
2. faça novo push em `main` com mudança no backend;
3. o workflow aplica update na mesma stack.

## 8. Troubleshooting

### Erro: ResourceNotFoundException em `aws lambda get-function`

A função em `FUNCTION_NAME` não existe na região configurada. Crie a função
uma vez (manual/CLI) e rode o workflow novamente.

### Erro: AccessDenied no deploy

Revise policy do usuário IAM do GitHub Actions e valide região/conta/ARNs.

### Erro: No updates are to be performed

Comportamento normal quando não há diferença no template/parâmetros/código
referenciado.

### Erro de CORS no frontend

Revise `ALLOWED_ORIGIN` (origem exata, sem `/` final) e faça novo deploy.

## 9. Segurança recomendada (próximo passo)

Para eliminar chaves estáticas no GitHub, migre para OIDC com
`aws-actions/configure-aws-credentials` assumindo role.
Este repositório já está estruturado de forma que essa troca é simples.

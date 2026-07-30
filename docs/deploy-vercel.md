# Deploy do frontend na Vercel

Este guia publica somente o frontend. Antes de começar, o backend deve estar
publicado e a URL da API deve estar disponível.

## Configuração

1. Envie este repositório para o GitHub, GitLab ou Bitbucket.
2. Na [Vercel](https://vercel.com), selecione **Add New > Project** e importe o
   repositório.
3. Configure o projeto:

   | Campo | Valor |
   | --- | --- |
   | Root Directory | `frontend` |
   | Framework Preset | `Vite` |
   | Build Command | `npm run build` |
   | Output Directory | `dist` |

4. Em **Environment Variables**, adicione:

   ```text
   VITE_API_BASE_URL=https://SUA-API.execute-api.REGIAO.amazonaws.com
   ```

   Use a URL sem `/todos` e sem `/` no final.

5. Clique em **Deploy**.

## Liberar o acesso no backend

Depois do deploy, copie a URL gerada pela Vercel e defina esse endereço na
variável `ALLOWED_ORIGIN` da Lambda:

```text
ALLOWED_ORIGIN=https://seu-projeto.vercel.app
```

Use somente a origem, sem caminho e sem `/` no final. Depois de salvar a
variável, abra o endereço da Vercel e teste a criação de uma tarefa.

## Atualizações

Novos commits enviados para a branch conectada geram novos deploys
automaticamente. Se alterar `VITE_API_BASE_URL`, faça um novo deploy em
**Deployments > Redeploy**, pois a variável é aplicada durante o build.

## Problemas comuns

- **Erro de CORS:** confira se `ALLOWED_ORIGIN` contém exatamente a URL da
  Vercel.
- **Frontend tenta acessar localhost:** confira `VITE_API_BASE_URL` e faça um
  novo deploy.
- **Erro 404 ou 5xx da API:** teste a URL do backend e confirme que ela não
  contém `/todos` no final.

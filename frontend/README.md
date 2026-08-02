# Frontend

Aplicação web Vite para a API de tarefas.

## Executar localmente

Inicie o backend primeiro. Depois, execute os comandos abaixo neste diretório:

```bash
npm install
npm run dev
```

Acesse `http://localhost:5173`. Por padrão, o frontend utiliza
`http://localhost:8081` como backend.

## Build de produção

```bash
npm run build
npm run preview
```

Os arquivos gerados são salvos no diretório `dist/`.

## Deploy na Vercel

Siga o [guia de deploy na Vercel](./docs/deploy-vercel.md).

## Deploy com Docker

```bash
docker build -t todos-frontend .
docker run --rm -p 8080:80 -e API_BASE_URL="https://sua-api.exemplo.com" todos-frontend
```

A imagem Docker lê `API_BASE_URL` quando o contêiner é iniciado. Na Vercel,
utilize `VITE_API_BASE_URL`, que é aplicada durante o build.

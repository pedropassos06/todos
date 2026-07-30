# Frontend

Vite web application for the Todos API.

## Run locally

Start the backend first. Then, from this directory:

```bash
npm install
npm run dev
```

Open `http://localhost:5173`. The frontend uses
`http://localhost:8081` by default, so no frontend environment file is
required for local development.

## Production build

```bash
npm run build
npm run preview
```

The generated files are written to `dist/`.

## Deploy to Vercel

Import the repository and configure:

- Root Directory: `frontend`
- Framework Preset: Vite
- Build Command: `npm run build`
- Output Directory: `dist`
- Environment Variable:
  `VITE_API_BASE_URL=https://SUA-API.execute-api.REGIAO.amazonaws.com`

After deployment, configure the Vercel URL as `ALLOWED_ORIGIN` in the
backend Lambda.

## Docker deployment

```bash
docker build -t todos-frontend .
docker run --rm -p 8080:80 -e API_BASE_URL="https://sua-api.exemplo.com" todos-frontend
```

The Docker image reads `API_BASE_URL` when its container starts. Vercel uses
`VITE_API_BASE_URL` during the build instead.

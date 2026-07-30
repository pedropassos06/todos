# Todos

Todo application with two independent components:

- `backend/`: Go API using DynamoDB.
- `frontend/`: Vite web application.

## Requirements

- Docker and Docker Compose
- Node.js 20+
- npm

Go is only required when building or testing the backend outside Docker.

## Run locally

### 1. Start the backend

Open a terminal from the project root:

```bash
cp .env.example .env
cd backend
make start
```

This starts the backend and DynamoDB through Docker Compose. The API is
available at `http://localhost:8081`.

The environment file only needs to be copied the first time.

### 2. Start the frontend

Open another terminal from the project root:

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:5173`.

`npm install` only needs to be run the first time or when dependencies change.
The frontend uses `http://localhost:8081` by default.

## Stop the backend

From `backend/`:

```bash
make stop
```

## Useful commands

Backend commands are run from `backend/`:

```bash
make start
make stop
make logs
make test
make package
```

Frontend commands are run from `frontend/`:

```bash
npm run dev
npm run build
npm run preview
```

## Production deployment

- Backend: AWS Lambda, API Gateway and DynamoDB.
- Frontend: Vercel.

See [`docs/deploy-manual-aws-console.md`](docs/deploy-manual-aws-console.md)
for the AWS instructions. The frontend deployment settings are documented in
[`frontend/README.md`](frontend/README.md).

# Backend

Go REST API backed by DynamoDB.

## Run locally

From the project root:

```bash
cp .env.example .env
cd backend
make start
```

The API is available at `http://localhost:8081`.

The environment file only needs to be copied the first time. `make start`
starts both the API and LocalStack DynamoDB.

## Commands

Run these commands from this directory:

```bash
make start    # Start the API and DynamoDB
make stop     # Stop the local services
make logs     # Follow API logs
make test     # Run Go tests
make build    # Build the AWS Lambda executable
make package  # Generate function.zip for AWS Lambda
make clean    # Remove backend build artifacts
```

The Lambda uses `TABLE_NAME` and `ALLOWED_ORIGIN` in production. Local values
are defined in the root `.env` file.

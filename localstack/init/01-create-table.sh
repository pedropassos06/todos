#!/bin/sh
set -e

TABLE_NAME="${TABLE_NAME:-todos-table}"

awslocal dynamodb describe-table --table-name "$TABLE_NAME" >/dev/null 2>&1 && exit 0

awslocal dynamodb create-table \
  --table-name "$TABLE_NAME" \
  --attribute-definitions AttributeName=id,AttributeType=S \
  --key-schema AttributeName=id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST >/dev/null

echo "Created DynamoDB table: $TABLE_NAME"

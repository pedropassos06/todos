package main

import (
	"context"
	"log"
	"os"

	"todos/internal/handler"
	"todos/internal/repository"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// commit
func main() {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Fatal("TABLE_NAME environment variable is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	dynamodbEndpoint := os.Getenv("DYNAMODB_ENDPOINT")
	dbClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if dynamodbEndpoint != "" {
			o.BaseEndpoint = aws.String(dynamodbEndpoint)
		}
	})
	repo := repository.NewDynamoDBRepository(dbClient, tableName)
	h := handler.New(repo)

	lambda.Start(h.HandleRequest)
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"todos/internal/handler"
	"todos/internal/repository"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body := ""
		if r.Body != nil {
			b, _ := io.ReadAll(r.Body)
			body = string(b)
		}

		path := r.URL.Path
		pathParams := map[string]string{}

		// detect /todos/{id}
		routePath := path
		if strings.HasPrefix(path, "/todos/") {
			id := strings.TrimPrefix(path, "/todos/")
			if id != "" {
				pathParams["id"] = id
				routePath = "/todos/{id}"
			}
		}

		routeKey := fmt.Sprintf("%s %s", r.Method, routePath)

		req := events.APIGatewayV2HTTPRequest{
			RawPath: path,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				RouteKey: routeKey,
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: r.Method,
					Path:   path,
				},
			},
			PathParameters: pathParams,
			Body:           body,
		}

		resp, err := h.HandleRequest(r.Context(), req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		fmt.Fprint(w, resp.Body)
	})

	log.Printf("server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

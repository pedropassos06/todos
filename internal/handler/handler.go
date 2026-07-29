package handler

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"todos/internal/model"
	"todos/internal/repository"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	repo *repository.DynamoDBRepository
}

type createTodoRequest struct {
	Title string `json:"title"`
}

type updateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

type todosResponse struct {
	Todos []model.Todo `json:"todos"`
}

func New(repo *repository.DynamoDBRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := req.RequestContext.HTTP.Method
	path := req.RawPath
	if path == "" {
		path = req.RequestContext.HTTP.Path
	}
	routeKey := req.RequestContext.RouteKey

	switch routeKey {
	case "GET /todos":
		return h.listTodos(ctx)
	case "POST /todos":
		return h.createTodo(ctx, req.Body)
	case "PATCH /todos/{id}":
		return h.updateTodo(ctx, req.PathParameters["id"], req.Body)
	case "DELETE /todos/{id}":
		return h.deleteTodo(ctx, req.PathParameters["id"])
	}

	// Fallback by method + path when routeKey is unavailable in local tests.
	switch {
	case method == http.MethodGet && path == "/todos":
		return h.listTodos(ctx)
	case method == http.MethodPost && path == "/todos":
		return h.createTodo(ctx, req.Body)
	case method == http.MethodPatch && strings.HasPrefix(path, "/todos/"):
		id := req.PathParameters["id"]
		if id == "" {
			id = strings.TrimPrefix(path, "/todos/")
		}
		return h.updateTodo(ctx, id, req.Body)
	case method == http.MethodDelete && strings.HasPrefix(path, "/todos/"):
		id := req.PathParameters["id"]
		if id == "" {
			id = strings.TrimPrefix(path, "/todos/")
		}
		return h.deleteTodo(ctx, id)
	default:
		return jsonResponse(http.StatusNotFound, map[string]string{"message": "not found"})
	}
}

func (h *Handler) listTodos(ctx context.Context) (events.APIGatewayV2HTTPResponse, error) {
	todos, err := h.repo.ListTodos(ctx)
	if err != nil {
		log.Printf("error listing todos: %v", err)
		return jsonResponse(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return jsonResponse(http.StatusOK, todosResponse{Todos: todos})
}

func (h *Handler) createTodo(ctx context.Context, body string) (events.APIGatewayV2HTTPResponse, error) {
	var input createTodoRequest
	if err := json.Unmarshal([]byte(body), &input); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return jsonResponse(http.StatusBadRequest, map[string]string{"message": "title is required"})
	}

	now := time.Now().UTC().Format(time.RFC3339)
	todo := model.Todo{
		ID:        newUUID(),
		Title:     title,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.repo.CreateTodo(ctx, todo); err != nil {
		log.Printf("error creating todo: %v", err)
		return jsonResponse(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return jsonResponse(http.StatusCreated, todo)
}

func (h *Handler) updateTodo(ctx context.Context, id string, body string) (events.APIGatewayV2HTTPResponse, error) {
	if strings.TrimSpace(id) == "" {
		return jsonResponse(http.StatusNotFound, map[string]string{"message": "not found"})
	}

	var input updateTodoRequest
	if err := json.Unmarshal([]byte(body), &input); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	if input.Title != nil {
		trimmed := strings.TrimSpace(*input.Title)
		if trimmed == "" {
			return jsonResponse(http.StatusBadRequest, map[string]string{"message": "title is required"})
		}
		input.Title = &trimmed
	}

	now := time.Now().UTC().Format(time.RFC3339)
	updatedTodo, err := h.repo.UpdateTodo(ctx, id, input.Title, input.Completed, now)
	if err != nil {
		if err == repository.ErrTodoNotFound {
			return jsonResponse(http.StatusNotFound, map[string]string{"message": "todo not found"})
		}
		log.Printf("error updating todo: %v", err)
		return jsonResponse(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return jsonResponse(http.StatusOK, updatedTodo)
}

func (h *Handler) deleteTodo(ctx context.Context, id string) (events.APIGatewayV2HTTPResponse, error) {
	if strings.TrimSpace(id) == "" {
		return jsonResponse(http.StatusNotFound, map[string]string{"message": "not found"})
	}

	if err := h.repo.DeleteTodo(ctx, id); err != nil {
		if err == repository.ErrTodoNotFound {
			return jsonResponse(http.StatusNotFound, map[string]string{"message": "todo not found"})
		}
		log.Printf("error deleting todo: %v", err)
		return jsonResponse(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNoContent,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func jsonResponse(status int, body any) (events.APIGatewayV2HTTPResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		log.Printf("error marshaling response: %v", err)
		fallback := []byte(`{"message":"internal server error"}`)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: string(fallback),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(payload),
	}, nil
}

func newUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate UUID")
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	)
}

package repository

import (
	"context"
	"errors"

	"todos/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ErrTodoNotFound = errors.New("todo not found")

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(client *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBRepository) ListTodos(ctx context.Context) ([]model.Todo, error) {
	out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, err
	}

	if len(out.Items) == 0 {
		return []model.Todo{}, nil
	}

	var todos []model.Todo
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *DynamoDBRepository) CreateTodo(ctx context.Context, todo model.Todo) error {
	item, err := attributevalue.MarshalMap(todo)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *DynamoDBRepository) GetTodo(ctx context.Context, id string) (*model.Todo, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(out.Item) == 0 {
		return nil, ErrTodoNotFound
	}

	var todo model.Todo
	if err := attributevalue.UnmarshalMap(out.Item, &todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *DynamoDBRepository) UpdateTodo(ctx context.Context, id string, title *string, completed *bool, updatedAt string) (*model.Todo, error) {
	expressionValues := map[string]types.AttributeValue{
		":updatedAt": &types.AttributeValueMemberS{Value: updatedAt},
	}

	updateExpression := "SET updatedAt = :updatedAt"

	if title != nil {
		expressionValues[":title"] = &types.AttributeValueMemberS{Value: *title}
		updateExpression += ", title = :title"
	}

	if completed != nil {
		expressionValues[":completed"] = &types.AttributeValueMemberBOOL{Value: *completed}
		updateExpression += ", completed = :completed"
	}

	out, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression:       aws.String("attribute_exists(id)"),
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		var conditionalCheckFailed *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailed) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	var todo model.Todo
	if err := attributevalue.UnmarshalMap(out.Attributes, &todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *DynamoDBRepository) DeleteTodo(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	})
	if err != nil {
		var conditionalCheckFailed *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailed) {
			return ErrTodoNotFound
		}
		return err
	}

	return nil
}

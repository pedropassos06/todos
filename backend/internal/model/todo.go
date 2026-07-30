package model

type Todo struct {
	ID        string `json:"id" dynamodbav:"id"`
	Title     string `json:"title" dynamodbav:"title"`
	Completed bool   `json:"completed" dynamodbav:"completed"`
	CreatedAt string `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt string `json:"updatedAt" dynamodbav:"updatedAt"`
}

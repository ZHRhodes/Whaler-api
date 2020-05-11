// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/heroku/whaler-api/models"
)

type NewTodo struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type Todo struct {
	ID   string       `json:"id"`
	Text string       `json:"text"`
	Done bool         `json:"done"`
	User *models.User `json:"user"`
}

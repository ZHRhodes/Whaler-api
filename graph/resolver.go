package graph

import (
	"github.com/heroku/whaler-api/graph/model"
	"github.com/jinzhu/gorm"
)

// This file will not be regenerated automatically.

// Resolver serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	db    *gorm.DB
	todos []*model.Todo
}

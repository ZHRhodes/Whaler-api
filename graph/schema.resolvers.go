package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/heroku/whaler-api/graph/generated"
	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/middleware"
	"github.com/heroku/whaler-api/models"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*models.User, error) {
	user, err := models.CreateUser(input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mutationResolver) CreateAccount(ctx context.Context, input model.NewAccount) (*models.Account, error) {
	account, err := models.CreateAccount(input)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *mutationResolver) CreateContact(ctx context.Context, input model.NewContact) (*models.Contact, error) {
	contact, err := models.CreateContact(input)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (r *mutationResolver) CreateWorkspace(ctx context.Context, input model.NewWorkspace) (*models.Workspace, error) {
	workspace, err := models.CreateWorkspace(input)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (r *queryResolver) Workspaces(ctx context.Context) ([]*models.Workspace, error) {
	userID := middleware.UserIDFromContext(ctx)
	return models.FetchWorkspaces(r.DB, userID)
}

func (r *queryResolver) Accounts(ctx context.Context) ([]*models.Account, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Organization(ctx context.Context) (*models.Organization, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/heroku/whaler-api/graph/generated"
	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/middleware"
	"github.com/heroku/whaler-api/models"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*models.User, error) {
	user, err := models.CreateUser(input.Email, input.Password, input.OrganizationID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mutationResolver) CreateOrganization(ctx context.Context, input model.NewOrganization) (*models.Organization, error) {
	return models.CreateOrganization(r.DB, input)
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

func (r *mutationResolver) CreateContactAssignmentEntry(ctx context.Context, input model.NewContactAssignmentEntry) (*models.ContactAssignmentEntry, error) {
	entry, err := models.CreateContactAssignmentEntry(input)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (r *mutationResolver) CreateAccountAssignmentEntry(ctx context.Context, input model.NewAccountAssignmentEntry) (*models.AccountAssignmentEntry, error) {
	return models.CreateAccountAssignmentEntry(input)
}

func (r *mutationResolver) SaveAccounts(ctx context.Context, input []*model.NewAccount) ([]*models.Account, error) {
	userID := middleware.UserIDFromContext(ctx)
	accounts, err := models.SaveAccounts(input, userID)
	if err != nil {
		return nil, err
	}
	return accounts, err
}

func (r *mutationResolver) SaveContacts(ctx context.Context, input []*model.NewContact) ([]*models.Contact, error) {
	return models.SaveContacts(input)
}

func (r *queryResolver) Workspaces(ctx context.Context) ([]*models.Workspace, error) {
	userID := middleware.UserIDFromContext(ctx)
	preloads := getPreloads(ctx)
	return models.FetchWorkspaces(r.DB, preloads, userID)
}

func (r *queryResolver) Organization(ctx context.Context) (*models.Organization, error) {
	userID := middleware.UserIDFromContext(ctx)
	user := models.FetchUser(userID)
	preloads := getPreloads(ctx)
	return models.FetchOrganization(r.DB, preloads, user.OrganizationID)
}

func (r *queryResolver) Accounts(ctx context.Context) ([]*models.Account, error) {
	userID := middleware.UserIDFromContext(ctx)
	return models.FetchAccounts(userID)
}

func (r *queryResolver) ContactAssignmentEntries(ctx context.Context, contactID string) ([]*models.ContactAssignmentEntry, error) {
	return models.FetchContactAssignmentEntries(contactID)
}

func (r *queryResolver) AccountAssignmentEntries(ctx context.Context, accountID string) ([]*models.AccountAssignmentEntry, error) {
	return models.FetchAccountAssignmentEntries(accountID)
}

// Account returns generated.AccountResolver implementation.
func (r *Resolver) Account() generated.AccountResolver { return &accountResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type accountResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func getPreloads(ctx context.Context) []string {
	return getNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}
func getNestedPreloads(ctx *graphql.RequestContext, fields []graphql.CollectedField, prefix string) (preloads []string) {
	for _, column := range fields {
		prefixColumn := getPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(preloads, getNestedPreloads(ctx, graphql.CollectFields(ctx, column.SelectionSet, nil), prefixColumn)...)
		preloads = append(preloads, getNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...)

	}
	return
}
func getPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}

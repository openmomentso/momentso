package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/openmomentso/momentso/pkg/app/auth"
	"github.com/openmomentso/momentso/pkg/database/db"
	"github.com/openmomentso/momentso/pkg/graph/model"
	"golang.org/x/crypto/bcrypt"
)

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, email string, password string) (*model.SignInPayload, error) {
	user, err := r.DB.UserFindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token := gonanoid.Must(64)
	_, err = r.DB.SessionCreate(ctx, db.SessionCreateParams{
		UserID: user.ID,
		Token:  token,
	})
	if err != nil {
		return nil, err
	}

	return &model.SignInPayload{
		User:  user,
		Token: token,
	}, nil
}

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, email string, password string) (*model.SignUpPayload, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := r.DB.UserCreate(ctx, db.UserCreateParams{
		Email:    email,
		Password: string(hash),
	})
	if err != nil {
		return nil, err
	}

	token := gonanoid.Must(64)
	_, err = r.DB.SessionCreate(ctx, db.SessionCreateParams{
		UserID: user.ID,
		Token:  token,
	})
	if err != nil {
		return nil, err
	}

	return &model.SignUpPayload{
		User:  user,
		Token: token,
	}, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*db.User, error) {
	user, ok := auth.UserForCtx(ctx)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	return &user, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

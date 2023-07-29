package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paw1a/grpc-media-converter/auth_service/internal/domain"
)

type Users interface {
	FindAll(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, userID int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	Delete(ctx context.Context, userID int64) error
}

type UserRepo struct {
	dbPool *pgxpool.Pool
}

func (u *UserRepo) FindAll(ctx context.Context) ([]domain.User, error) {
	return nil, nil
}

func (u *UserRepo) FindByID(ctx context.Context, userID int64) (domain.User, error) {
	return domain.User{}, nil
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return domain.User{}, nil
}

func (u *UserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	return domain.User{}, nil
}

func (u *UserRepo) Delete(ctx context.Context, userID int64) error {
	return nil
}

func NewUserRepository(dbPool *pgxpool.Pool) *UserRepo {
	return &UserRepo{dbPool: dbPool}
}

package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: nil,
	}
}

type CreateUserParams struct {
	FirebaseUid  string
	Email        string
	RevenueCatId *string
	Timezone     *string
	Locale       *string
}

func (repo *UserRepository) CreateUser(ctx context.Context, p *CreateUserParams) {
	sql, args, _ := sq.Insert("users").
		Columns("email, firebase_uid, revenuecat_id, timezone, locale").
		Values(p.Email, p.FirebaseUid, p.RevenueCatId, p.Timezone, p.Locale).
		Suffix("returning *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	err := repo.db.QueryRow(ctx, sql, args).Scan()
	if err != nil {
		fmt.Println(err)
		return
	}
}

type userDbModel struct {
	Id           string     `db:"id"`
	Email        string     `db:"email"`
	FirebaseUid  *string    `db:"firebase_uid"`
	RevenueCatId *string    `db:"revenuecat_id"`
	Timezone     *string    `db:"timezone"`
	Locale       *string    `db:"locale"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

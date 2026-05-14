package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/artivisi/tech-stack-2026/golang/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type Registration struct {
	db *sql.DB
}

func NewRegistration(db *sql.DB) *Registration {
	return &Registration{db: db}
}

func (r *Registration) Insert(ctx context.Context, reg domain.Registration) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO registration (id, email, full_name, phone, created_at) VALUES ($1, $2, $3, $4, $5)`,
		reg.ID, reg.Email, reg.FullName, reg.Phone, reg.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *Registration) FindAll(ctx context.Context) ([]domain.Registration, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, email, full_name, phone, created_at FROM registration ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		if err := rows.Scan(&reg.ID, &reg.Email, &reg.FullName, &reg.Phone, &reg.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, reg)
	}
	return results, rows.Err()
}

func (r *Registration) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return r.db.PingContext(ctx)
}

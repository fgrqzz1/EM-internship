package repository

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSubsriptionNotFound = errors.New("subsription not found")

type SubsriptionRepository struct {
	db *pgxpool.Pool
}

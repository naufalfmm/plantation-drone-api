// This file contains the repository implementation layer.
package repository

import (
	_ "github.com/lib/pq"
	"github.com/naufalfmm/plantation-drone-api/utils/db"
)

type Repository struct {
	Db db.DB
}

type NewRepositoryOptions struct {
	Db db.DB
}

func NewRepository(opts NewRepositoryOptions) *Repository {
	return &Repository{
		Db: opts.Db,
	}
}

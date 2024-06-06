package repository

import "context"

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateEstate(ctx context.Context, input CreateEstateInput) (err error) {
	err = r.Db.QueryRowContext(ctx, `INSERT INTO estates (id, width, length, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`, input.Id, input.Width, input.Length).Err()
	if err != nil {
		return
	}

	return
}

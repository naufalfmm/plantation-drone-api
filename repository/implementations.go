package repository

import (
	"context"
	"database/sql"
)

func (r *Repository) CreateEstate(ctx context.Context, input CreateEstateInput) (err error) {
	err = r.Db.QueryRowContext(ctx, `INSERT INTO estates (id, width, length, drone_distance, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())`, input.Id, input.Width, input.Length, ((input.Length-1)*10*input.Width + (input.Width-1)*10 + 2)).Err()
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetEstateById(ctx context.Context, input GetEstateByIdInput) (output GetEstateByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, `SELECT id, width, length, count, max, min, median, drone_distance FROM estates WHERE id = $1`, input.Id).Scan(&output.Id, &output.Width, &output.Length, &output.Count, &output.Max, &output.Min, &output.Median, &output.DroneDistance)
	if err != nil {
		return
	}

	return
}

func (r *Repository) CountCoordinateTree(ctx context.Context, input CountCoordinateTreeInput) (output CountCoordinateTreeOutput, err error) {
	err = r.Db.QueryRowContext(ctx, `SELECT COUNT(id) FROM estate_trees WHERE x = $1 AND y = $2`, input.X, input.Y).Scan(&output.Count)
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetPrevNextTree(ctx context.Context, input GetPrevNextTreeInput) (output GetPrevNextTreeOutput, err error) {
	stmts, err := r.Db.QueryContext(ctx, `SELECT x, y, height FROM estate_trees WHERE (x = $1 AND y = $2) OR (x = $3 AND y = $4) LIMIT 2`, input.PrevX, input.PrevY, input.NextX, input.NextY)
	if err != nil {
		return
	}

	for stmts.Next() {
		x, y, height := 0, 0, 0

		err = stmts.Scan(&x, &y, &height)
		if err != nil {
			return
		}

		if x == input.PrevX && y == input.PrevY {
			output.PrevTreeHeight = height
		}

		if x == input.NextX && y == input.NextY {
			output.NextTreeHeight = height
		}
	}

	return
}

func (r *Repository) CreateTree(ctx context.Context, input CreateTreeInput) (err error) {
	tx, err := r.Db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `UPDATE estates
		SET count = count + 1,
			max = CASE WHEN max < $1 THEN $1 ELSE max END,
			min = CASE WHEN (min = 0 OR min > $1) THEN $1 ELSE min END,
			drone_distance = drone_distance + $2,
			median = 0,
			updated_at = NOW()
		WHERE id = $3
	`, input.Height, input.DroneDistFactor, input.EstateId)
	if err != nil {
		return
	}
	rows.Close()

	rows, err = tx.QueryContext(ctx, `INSERT INTO estate_trees (id, estate_id, x, y, height, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`, input.Id, input.EstateId, input.X, input.Y, input.Height)
	if err != nil {
		return
	}
	rows.Close()

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetHeightEstateTrees(ctx context.Context, input GetHeightEstateTreesInput) (output GetHeightEstateTreesOutput, err error) {
	rows, err := r.Db.QueryContext(ctx, `SELECT height FROM estate_trees WHERE estate_id = $1`, input.EstateId)
	if err != nil {
		return
	}

	for rows.Next() {
		height := 0
		err = rows.Scan(&height)
		if err != nil {
			return
		}

		output.Heights = append(output.Heights, height)
	}

	return
}

func (r *Repository) StoreMedianEstate(ctx context.Context, input StoreMedianEstateInput) (err error) {
	err = r.Db.QueryRowContext(ctx, `UPDATE estates SET median = $1 WHERE id = $2`, input.Median, input.EstateId).Err()
	if err != nil {
		return
	}

	return
}

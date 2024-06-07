package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/naufalfmm/plantation-drone-api/utils/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateEstate(t *testing.T) {
	t.Run("Return no error when insert is success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRow := db.NewMockRow(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		input := CreateEstateInput{
			Id:     "aaaaa-bbbbb-ccccc-ddddd",
			Width:  5,
			Length: 6,
		}

		ctx := context.Background()

		mockDb.EXPECT().QueryRowContext(ctx, `INSERT INTO estates (id, width, length, drone_distance, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())`, input.Id, input.Width, input.Length, 292).Return(mockRow)
		mockRow.EXPECT().Err().Return(nil)

		err := repo.CreateEstate(ctx, input)

		assert.Nil(t, err)
	})

	t.Run("Return error when row error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRow := db.NewMockRow(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := CreateEstateInput{
			Id:     "aaaaa-bbbbb-ccccc-ddddd",
			Width:  5,
			Length: 6,
		}

		ctx := context.Background()

		mockDb.EXPECT().QueryRowContext(ctx, `INSERT INTO estates (id, width, length, drone_distance, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())`, input.Id, input.Width, input.Length, 292).Return(mockRow)
		mockRow.EXPECT().Err().Return(errAny)

		err := repo.CreateEstate(ctx, input)

		assert.Equal(t, errAny, err)
	})
}

func TestGetEstateById(t *testing.T) {
	t.Run("Return no error when get is success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRow := db.NewMockRow(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		input := GetEstateByIdInput{
			Id: "aaaaa-bbbbb-ccccc-ddddd",
		}

		expOutput := GetEstateByIdOutput{}

		ctx := context.Background()

		out := GetEstateByIdOutput{}
		mockDb.EXPECT().QueryRowContext(ctx, `SELECT id, width, length, count, max, min, median, drone_distance FROM estates WHERE id = $1`, input.Id).Return(mockRow)
		mockRow.EXPECT().Scan(&out.Id, &out.Width, &out.Length, &out.Count, &out.Max, &out.Min, &out.Median, &out.DroneDistance).Return(nil)

		output, err := repo.GetEstateById(ctx, input)

		assert.Nil(t, err)
		assert.Equal(t, expOutput, output)
	})

	t.Run("Return error when scan error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRow := db.NewMockRow(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := GetEstateByIdInput{
			Id: "aaaaa-bbbbb-ccccc-ddddd",
		}

		ctx := context.Background()

		out := GetEstateByIdOutput{}
		mockDb.EXPECT().QueryRowContext(ctx, `SELECT id, width, length, count, max, min, median, drone_distance FROM estates WHERE id = $1`, input.Id).Return(mockRow)
		mockRow.EXPECT().Scan(&out.Id, &out.Width, &out.Length, &out.Count, &out.Max, &out.Min, &out.Median, &out.DroneDistance).Return(errAny)

		output, err := repo.GetEstateById(ctx, input)

		assert.Equal(t, errAny, err)
		assert.Equal(t, GetEstateByIdOutput{}, output)
	})
}

func TestCountCoordinateTree(t *testing.T) {
	t.Run("Return no error when get is success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRow := db.NewMockRow(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		input := CountCoordinateTreeInput{
			X: 2,
			Y: 1,
		}

		ctx := context.Background()

		var count int
		mockDb.EXPECT().QueryRowContext(ctx, `SELECT COUNT(id) FROM estate_trees WHERE x = $1 AND y = $2`, input.X, input.Y).Return(mockRow)
		mockRow.EXPECT().Scan(&count).Return(nil)

		output, err := repo.CountCoordinateTree(ctx, input)

		assert.Nil(t, err)
		assert.Equal(t, CountCoordinateTreeOutput{}, output)
	})

	t.Run("Return error when scan error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRow := db.NewMockRow(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := CountCoordinateTreeInput{
			X: 2,
			Y: 1,
		}

		ctx := context.Background()

		var count int
		mockDb.EXPECT().QueryRowContext(ctx, `SELECT COUNT(id) FROM estate_trees WHERE x = $1 AND y = $2`, input.X, input.Y).Return(mockRow)
		mockRow.EXPECT().Scan(&count).Return(errAny)

		output, err := repo.CountCoordinateTree(ctx, input)

		assert.Equal(t, errAny, err)
		assert.Equal(t, CountCoordinateTreeOutput{}, output)
	})
}

func TestGetPrevNextTree(t *testing.T) {
	t.Run("Return the data when get is success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		input := GetPrevNextTreeInput{
			PrevX: 2,
			PrevY: 1,

			NextX: 4,
			NextY: 1,
		}

		expOutput := GetPrevNextTreeOutput{
			PrevTreeHeight: 10,
			NextTreeHeight: 15,
		}

		ctx := context.Background()

		var x, y, height int
		mockDb.EXPECT().QueryContext(ctx, `SELECT x, y, height FROM estate_trees WHERE (x = $1 AND y = $2) OR (x = $3 AND y = $4) LIMIT 2`, input.PrevX, input.PrevY, input.NextX, input.NextY).Return(mockRows, nil)
		mockRows.EXPECT().Next().Return(true)
		mockRows.EXPECT().Scan(&x, &y, &height).DoAndReturn(func(args ...interface{}) interface{} {
			*(args[0].(*int)) = input.PrevX
			*(args[1].(*int)) = input.PrevY
			*(args[2].(*int)) = expOutput.PrevTreeHeight

			return nil
		})
		mockRows.EXPECT().Next().Return(true)
		mockRows.EXPECT().Scan(&x, &y, &height).DoAndReturn(func(args ...interface{}) interface{} {
			*(args[0].(*int)) = input.NextX
			*(args[1].(*int)) = input.NextY
			*(args[2].(*int)) = expOutput.NextTreeHeight

			return nil
		})
		mockRows.EXPECT().Next().Return(false)

		output, err := repo.GetPrevNextTree(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expOutput, output)
	})

	t.Run("Return error when one of the scan is error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := GetPrevNextTreeInput{
			PrevX: 2,
			PrevY: 1,

			NextX: 4,
			NextY: 1,
		}

		expOutput := GetPrevNextTreeOutput{}

		ctx := context.Background()

		var x, y, height int
		mockDb.EXPECT().QueryContext(ctx, `SELECT x, y, height FROM estate_trees WHERE (x = $1 AND y = $2) OR (x = $3 AND y = $4) LIMIT 2`, input.PrevX, input.PrevY, input.NextX, input.NextY).Return(mockRows, nil)
		mockRows.EXPECT().Next().Return(true)
		mockRows.EXPECT().Scan(&x, &y, &height).Return(errAny)

		output, err := repo.GetPrevNextTree(ctx, input)

		assert.Equal(t, errAny, err)
		assert.Equal(t, expOutput, output)
	})

	t.Run("Return error when query context error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := GetPrevNextTreeInput{
			PrevX: 2,
			PrevY: 1,

			NextX: 4,
			NextY: 1,
		}

		expOutput := GetPrevNextTreeOutput{}

		ctx := context.Background()

		mockDb.EXPECT().QueryContext(ctx, `SELECT x, y, height FROM estate_trees WHERE (x = $1 AND y = $2) OR (x = $3 AND y = $4) LIMIT 2`, input.PrevX, input.PrevY, input.NextX, input.NextY).Return(mockRows, errAny)

		output, err := repo.GetPrevNextTree(ctx, input)

		assert.Equal(t, errAny, err)
		assert.Equal(t, expOutput, output)
	})
}

func TestCreateTree(t *testing.T) {
	t.Run("Return the tree when no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockTx := db.NewMockTx(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		input := CreateTreeInput{
			Id:     "aaaaa-bbbbb-ccccc-ddddd",
			X:      2,
			Y:      3,
			Height: 10,

			EstateId:        "bbbbb-ccccc-ddddd-eeeee",
			DroneDistFactor: 6,
		}

		ctx := context.Background()

		mockDb.EXPECT().BeginTx(ctx, &sql.TxOptions{}).Return(mockTx, nil)
		mockTx.EXPECT().Rollback()
		mockTx.EXPECT().QueryContext(ctx, `UPDATE estates
		SET count = count + 1,
			max = CASE WHEN max < $1 THEN $1 ELSE max END,
			min = CASE WHEN (min = 0 OR min > $1) THEN $1 ELSE min END,
			drone_distance = drone_distance + $2,
			median = 0,
			updated_at = NOW()
		WHERE id = $3
	`, input.Height, input.DroneDistFactor, input.EstateId).Return(mockRows, nil)
		mockRows.EXPECT().Close()
		mockTx.EXPECT().QueryContext(ctx, `INSERT INTO estate_trees (id, estate_id, x, y, height, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`, input.Id, input.EstateId, input.X, input.Y, input.Height).Return(mockRows, nil)
		mockRows.EXPECT().Close()
		mockTx.EXPECT().Commit().Return(nil)

		err := repo.CreateTree(ctx, input)

		assert.Nil(t, err)
	})

	t.Run("Return the tree when commit errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockTx := db.NewMockTx(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := CreateTreeInput{
			Id:     "aaaaa-bbbbb-ccccc-ddddd",
			X:      2,
			Y:      3,
			Height: 10,

			EstateId:        "bbbbb-ccccc-ddddd-eeeee",
			DroneDistFactor: 6,
		}

		ctx := context.Background()

		mockDb.EXPECT().BeginTx(ctx, &sql.TxOptions{}).Return(mockTx, nil)
		mockTx.EXPECT().Rollback()
		mockTx.EXPECT().QueryContext(ctx, `UPDATE estates
		SET count = count + 1,
			max = CASE WHEN max < $1 THEN $1 ELSE max END,
			min = CASE WHEN (min = 0 OR min > $1) THEN $1 ELSE min END,
			drone_distance = drone_distance + $2,
			median = 0,
			updated_at = NOW()
		WHERE id = $3
	`, input.Height, input.DroneDistFactor, input.EstateId).Return(mockRows, nil)
		mockRows.EXPECT().Close()
		mockTx.EXPECT().QueryContext(ctx, `INSERT INTO estate_trees (id, estate_id, x, y, height, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`, input.Id, input.EstateId, input.X, input.Y, input.Height).Return(mockRows, nil)
		mockRows.EXPECT().Close()
		mockTx.EXPECT().Commit().Return(errAny)

		err := repo.CreateTree(ctx, input)

		assert.Equal(t, errAny, err)
	})

	t.Run("Return the tree when query context of insert estate trees errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockTx := db.NewMockTx(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := CreateTreeInput{
			Id:     "aaaaa-bbbbb-ccccc-ddddd",
			X:      2,
			Y:      3,
			Height: 10,

			EstateId:        "bbbbb-ccccc-ddddd-eeeee",
			DroneDistFactor: 6,
		}

		ctx := context.Background()

		mockDb.EXPECT().BeginTx(ctx, &sql.TxOptions{}).Return(mockTx, nil)
		mockTx.EXPECT().Rollback()
		mockTx.EXPECT().QueryContext(ctx, `UPDATE estates
		SET count = count + 1,
			max = CASE WHEN max < $1 THEN $1 ELSE max END,
			min = CASE WHEN (min = 0 OR min > $1) THEN $1 ELSE min END,
			drone_distance = drone_distance + $2,
			median = 0,
			updated_at = NOW()
		WHERE id = $3
	`, input.Height, input.DroneDistFactor, input.EstateId).Return(mockRows, nil)
		mockRows.EXPECT().Close()
		mockTx.EXPECT().QueryContext(ctx, `INSERT INTO estate_trees (id, estate_id, x, y, height, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`, input.Id, input.EstateId, input.X, input.Y, input.Height).Return(mockRows, errAny)

		err := repo.CreateTree(ctx, input)

		assert.Equal(t, errAny, err)
	})

	t.Run("Return the tree when query context of update estates errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockTx := db.NewMockTx(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := CreateTreeInput{
			Id:     "aaaaa-bbbbb-ccccc-ddddd",
			X:      2,
			Y:      3,
			Height: 10,

			EstateId:        "bbbbb-ccccc-ddddd-eeeee",
			DroneDistFactor: 6,
		}

		ctx := context.Background()

		mockDb.EXPECT().BeginTx(ctx, &sql.TxOptions{}).Return(mockTx, nil)
		mockTx.EXPECT().Rollback()
		mockTx.EXPECT().QueryContext(ctx, `UPDATE estates
		SET count = count + 1,
			max = CASE WHEN max < $1 THEN $1 ELSE max END,
			min = CASE WHEN (min = 0 OR min > $1) THEN $1 ELSE min END,
			drone_distance = drone_distance + $2,
			median = 0,
			updated_at = NOW()
		WHERE id = $3
	`, input.Height, input.DroneDistFactor, input.EstateId).Return(mockRows, errAny)

		err := repo.CreateTree(ctx, input)

		assert.Equal(t, errAny, err)
	})

	t.Run("Return the tree when trx creating errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockTx := db.NewMockTx(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := CreateTreeInput{
			Id:     "aaaaa-bbbbb-ccccc-ddddd",
			X:      2,
			Y:      3,
			Height: 10,

			EstateId:        "bbbbb-ccccc-ddddd-eeeee",
			DroneDistFactor: 6,
		}

		ctx := context.Background()

		mockDb.EXPECT().BeginTx(ctx, &sql.TxOptions{}).Return(mockTx, errAny)

		err := repo.CreateTree(ctx, input)

		assert.Equal(t, errAny, err)
	})
}

func TestGetHeightEstateTrees(t *testing.T) {
	t.Run("Return the height when no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		input := GetHeightEstateTreesInput{
			EstateId: "bbbbb-ccccc-ddddd-eeeee",
		}

		expOutput := GetHeightEstateTreesOutput{
			Heights: []int{0},
		}

		ctx := context.Background()

		var height int
		mockDb.EXPECT().QueryContext(ctx, `SELECT height FROM estate_trees WHERE estate_id = $1`, input.EstateId).Return(mockRows, nil)
		mockRows.EXPECT().Next().Return(true)
		mockRows.EXPECT().Scan(&height).Return(nil)
		mockRows.EXPECT().Next().Return(false)

		output, err := repo.GetHeightEstateTrees(ctx, input)

		assert.Nil(t, err)
		assert.Equal(t, expOutput, output)
	})

	t.Run("Return error when scan errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := GetHeightEstateTreesInput{
			EstateId: "bbbbb-ccccc-ddddd-eeeee",
		}

		expOutput := GetHeightEstateTreesOutput{}

		ctx := context.Background()

		var height int
		mockDb.EXPECT().QueryContext(ctx, `SELECT height FROM estate_trees WHERE estate_id = $1`, input.EstateId).Return(mockRows, nil)
		mockRows.EXPECT().Next().Return(true)
		mockRows.EXPECT().Scan(&height).Return(errAny)

		output, err := repo.GetHeightEstateTrees(ctx, input)

		assert.Equal(t, errAny, err)
		assert.Equal(t, expOutput, output)
	})

	t.Run("Return error when query context errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRows := db.NewMockRows(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := GetHeightEstateTreesInput{
			EstateId: "bbbbb-ccccc-ddddd-eeeee",
		}

		expOutput := GetHeightEstateTreesOutput{}

		ctx := context.Background()

		mockDb.EXPECT().QueryContext(ctx, `SELECT height FROM estate_trees WHERE estate_id = $1`, input.EstateId).Return(mockRows, errAny)

		output, err := repo.GetHeightEstateTrees(ctx, input)

		assert.Equal(t, errAny, err)
		assert.Equal(t, expOutput, output)
	})
}

func TestStoreMedianEstate(t *testing.T) {
	t.Run("Return no error when update is success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRow := db.NewMockRow(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		input := StoreMedianEstateInput{
			EstateId: "aaaaa-bbbbb-ccccc-ddddd",
			Median:   5.5,
		}

		ctx := context.Background()

		mockDb.EXPECT().QueryRowContext(ctx, `UPDATE estates SET median = $1 WHERE id = $2`, input.Median, input.EstateId).Return(mockRow)
		mockRow.EXPECT().Err().Return(nil)

		err := repo.StoreMedianEstate(ctx, input)

		assert.Nil(t, err)
	})

	t.Run("Return error when row error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)
		mockRow := db.NewMockRow(ctrl)

		repo := Repository{
			Db: mockDb,
		}

		errAny := errors.New("any error")

		input := StoreMedianEstateInput{
			EstateId: "aaaaa-bbbbb-ccccc-ddddd",
			Median:   5.5,
		}

		ctx := context.Background()

		mockDb.EXPECT().QueryRowContext(ctx, `UPDATE estates SET median = $1 WHERE id = $2`, input.Median, input.EstateId).Return(mockRow)
		mockRow.EXPECT().Err().Return(errAny)

		err := repo.StoreMedianEstate(ctx, input)

		assert.Equal(t, errAny, err)
	})
}

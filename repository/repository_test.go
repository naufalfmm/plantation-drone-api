// This file contains the repository implementation layer.
package repository

import (
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"github.com/naufalfmm/plantation-drone-api/utils/db"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	t.Run("Return repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDb := db.NewMockDB(ctrl)

		expRepo := &Repository{
			Db: mockDb,
		}

		repo := NewRepository(NewRepositoryOptions{
			Db: mockDb,
		})

		assert.Equal(t, expRepo, repo)
	})
}

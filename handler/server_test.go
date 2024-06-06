package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/naufalfmm/plantation-drone-api/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("Create New Server", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := repository.NewMockRepositoryInterface(ctrl)

		expServer := &Server{
			Repository: repoMock,
		}

		server := NewServer(NewServerOptions{
			Repository: repoMock,
		})

		assert.Equal(t, expServer, server)
	})
}

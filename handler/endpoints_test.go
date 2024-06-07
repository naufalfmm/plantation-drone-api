package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/naufalfmm/plantation-drone-api/generated"
	"github.com/naufalfmm/plantation-drone-api/repository"
	"github.com/stretchr/testify/assert"
)

func readJsonResult(t *testing.T, resp *http.Response) map[string]any {
	var result map[string]any
	err := json.NewDecoder(resp.Body).Decode(&result)
	defer resp.Body.Close()

	assert.NoError(t, err)

	return result
}

func readJson[T any](t *testing.T, resp *http.Response) T {
	var result T
	err := json.NewDecoder(resp.Body).Decode(&result)
	defer resp.Body.Close()

	assert.NoError(t, err)

	return result
}

func TestGetHello(t *testing.T) {
	t.Run("Return 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		params := generated.GetHelloParams{
			Id: 15,
		}

		expResp := generated.HelloResponse{
			Message: fmt.Sprintf("Hello User %d", params.Id),
		}

		err := server.GetHello(ec, params)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resRecorder.Code)
		assert.Equal(t, expResp.Message, resp["message"])
	})
}

func TestPostEstate(t *testing.T) {
	t.Run("Return 201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"length\": 6, \"width\": 6}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		mockRepo.EXPECT().CreateEstate(ec.Request().Context(), gomock.Any()).Return(nil)

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resRecorder.Code)
		assert.NotEmpty(t, resp["id"])

		_, err = uuid.Parse(resp["id"].(string))
		assert.Nil(t, err)
	})

	t.Run("Return 500 when create estate error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"length\": 6, \"width\": 6}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		anyErr := errors.New("any error")

		mockRepo.EXPECT().CreateEstate(ec.Request().Context(), gomock.Any()).Return(anyErr)

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resRecorder.Code)
		assert.Equal(t, anyErr.Error(), resp["message"])
	})

	t.Run("Return 400 when width is negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"length\": 6, \"width\": -6}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrNegativeZeroBuilder("width").Error(), resp["message"])
	})

	t.Run("Return 400 when width is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"length\": 6, \"width\": 0}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrNegativeZeroBuilder("width").Error(), resp["message"])
	})

	t.Run("Return 400 when width is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"length\": 6}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrNegativeZeroBuilder("width").Error(), resp["message"])
	})

	t.Run("Return 400 when length is negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"length\": -6, \"width\": 6}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrNegativeZeroBuilder("length").Error(), resp["message"])
	})

	t.Run("Return 400 when length is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"length\": 0, \"width\": 6}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrNegativeZeroBuilder("length").Error(), resp["message"])
	})

	t.Run("Return 400 when length is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"width\": 6}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrNegativeZeroBuilder("length").Error(), resp["message"])
	})

	t.Run("Return 400 when header missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"length\": 6, \"width\": 6}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		err := server.PostEstate(ec)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, "code=415, message=Unsupported Media Type", resp["message"])
	})
}

func TestPostEstateIdTree(t *testing.T) {
	t.Run("Return 201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      2,
			Y:      1,
			Height: 11,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 0,
		}, nil)
		mockRepo.EXPECT().GetPrevNextTree(ec.Request().Context(), repository.GetPrevNextTreeInput{
			PrevX: 1,
			PrevY: 1,
			NextX: 3,
			NextY: 1,
		}).Return(repository.GetPrevNextTreeOutput{
			PrevTreeHeight: 0,
			NextTreeHeight: 0,
		}, nil)
		mockRepo.EXPECT().CreateTree(ec.Request().Context(), gomock.Any()).Return(nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resRecorder.Code)
		assert.NotEmpty(t, resp["id"])

		_, err = uuid.Parse(resp["id"].(string))
		assert.Nil(t, err)
	})

	t.Run("Return 201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 1, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      1,
			Y:      1,
			Height: 11,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 0,
		}, nil)
		mockRepo.EXPECT().GetPrevNextTree(ec.Request().Context(), repository.GetPrevNextTreeInput{
			PrevX: 1,
			PrevY: 0,
			NextX: 2,
			NextY: 1,
		}).Return(repository.GetPrevNextTreeOutput{
			PrevTreeHeight: 0,
			NextTreeHeight: 0,
		}, nil)
		mockRepo.EXPECT().CreateTree(ec.Request().Context(), gomock.Any()).Return(nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resRecorder.Code)
		assert.NotEmpty(t, resp["id"])

		_, err = uuid.Parse(resp["id"].(string))
		assert.Nil(t, err)
	})

	t.Run("Return 201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 6, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      6,
			Y:      1,
			Height: 11,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 0,
		}, nil)
		mockRepo.EXPECT().GetPrevNextTree(ec.Request().Context(), repository.GetPrevNextTreeInput{
			PrevX: 5,
			PrevY: 1,
			NextX: 6,
			NextY: 2,
		}).Return(repository.GetPrevNextTreeOutput{
			PrevTreeHeight: 0,
			NextTreeHeight: 0,
		}, nil)
		mockRepo.EXPECT().CreateTree(ec.Request().Context(), gomock.Any()).Return(nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resRecorder.Code)
		assert.NotEmpty(t, resp["id"])

		_, err = uuid.Parse(resp["id"].(string))
		assert.Nil(t, err)
	})

	t.Run("Return 201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 6, \"y\": 2, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      6,
			Y:      2,
			Height: 11,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 0,
		}, nil)
		mockRepo.EXPECT().GetPrevNextTree(ec.Request().Context(), repository.GetPrevNextTreeInput{
			PrevX: 6,
			PrevY: 1,
			NextX: 5,
			NextY: 2,
		}).Return(repository.GetPrevNextTreeOutput{
			PrevTreeHeight: 0,
			NextTreeHeight: 0,
		}, nil)
		mockRepo.EXPECT().CreateTree(ec.Request().Context(), gomock.Any()).Return(nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resRecorder.Code)
		assert.NotEmpty(t, resp["id"])

		_, err = uuid.Parse(resp["id"].(string))
		assert.Nil(t, err)
	})

	t.Run("Return 201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 3, \"y\": 2, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      3,
			Y:      2,
			Height: 11,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 0,
		}, nil)
		mockRepo.EXPECT().GetPrevNextTree(ec.Request().Context(), repository.GetPrevNextTreeInput{
			PrevX: 4,
			PrevY: 2,
			NextX: 2,
			NextY: 2,
		}).Return(repository.GetPrevNextTreeOutput{
			PrevTreeHeight: 0,
			NextTreeHeight: 0,
		}, nil)
		mockRepo.EXPECT().CreateTree(ec.Request().Context(), gomock.Any()).Return(nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resRecorder.Code)
		assert.NotEmpty(t, resp["id"])

		_, err = uuid.Parse(resp["id"].(string))
		assert.Nil(t, err)
	})

	t.Run("Return 201", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 1, \"y\": 2, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      1,
			Y:      2,
			Height: 11,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 0,
		}, nil)
		mockRepo.EXPECT().GetPrevNextTree(ec.Request().Context(), repository.GetPrevNextTreeInput{
			PrevX: 2,
			PrevY: 2,
			NextX: 1,
			NextY: 3,
		}).Return(repository.GetPrevNextTreeOutput{
			PrevTreeHeight: 0,
			NextTreeHeight: 0,
		}, nil)
		mockRepo.EXPECT().CreateTree(ec.Request().Context(), gomock.Any()).Return(nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resRecorder.Code)
		assert.NotEmpty(t, resp["id"])

		_, err = uuid.Parse(resp["id"].(string))
		assert.Nil(t, err)
	})

	t.Run("Return 500 when create tree error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      2,
			Y:      1,
			Height: 11,
		}

		id := uuid.New().String()

		anyErr := errors.New("any error")

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 0,
		}, nil)
		mockRepo.EXPECT().GetPrevNextTree(ec.Request().Context(), repository.GetPrevNextTreeInput{
			PrevX: 1,
			PrevY: 1,
			NextX: 3,
			NextY: 1,
		}).Return(repository.GetPrevNextTreeOutput{
			PrevTreeHeight: 0,
			NextTreeHeight: 0,
		}, nil)
		mockRepo.EXPECT().CreateTree(ec.Request().Context(), gomock.Any()).Return(anyErr)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resRecorder.Code)
		assert.Equal(t, anyErr.Error(), resp["message"])
	})

	t.Run("Return 500 when get prev next error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      2,
			Y:      1,
			Height: 11,
		}

		id := uuid.New().String()

		anyErr := errors.New("any error")

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 0,
		}, nil)
		mockRepo.EXPECT().GetPrevNextTree(ec.Request().Context(), repository.GetPrevNextTreeInput{
			PrevX: 1,
			PrevY: 1,
			NextX: 3,
			NextY: 1,
		}).Return(repository.GetPrevNextTreeOutput{}, anyErr)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resRecorder.Code)
		assert.Equal(t, anyErr.Error(), resp["message"])
	})

	t.Run("Return 400 when tree exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      2,
			Y:      1,
			Height: 11,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{
			Count: 1,
		}, nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrTreeExist.Error(), resp["message"])
	})

	t.Run("Return 500 when count tree error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		bodyReq := generated.CreateTreeRequest{
			X:      2,
			Y:      1,
			Height: 11,
		}

		id := uuid.New().String()

		anyErr := errors.New("any error")

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)
		mockRepo.EXPECT().CountCoordinateTree(ec.Request().Context(), repository.CountCoordinateTreeInput{
			X: bodyReq.X,
			Y: bodyReq.Y,
		}).Return(repository.CountCoordinateTreeOutput{}, anyErr)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resRecorder.Code)
		assert.Equal(t, anyErr.Error(), resp["message"])
	})

	t.Run("Return 400 when x out of bound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 7, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrCoordinateOutOfBound.Error(), resp["message"])
	})

	t.Run("Return 400 when y out of bound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 7, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{
			Length: 6,
			Width:  6,
		}, nil)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrCoordinateOutOfBound.Error(), resp["message"])
	})

	t.Run("Return 500 when get estate error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		errAny := errors.New("any error")

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{}, errAny)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resRecorder.Code)
		assert.Equal(t, errAny.Error(), resp["message"])
	})

	t.Run("Return 400 when estate missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(repository.GetEstateByIdOutput{}, sql.ErrNoRows)

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resRecorder.Code)
		assert.Equal(t, ErrNotFoundBuilder("estate").Error(), resp["message"])
	})

	t.Run("Return 400 when height is greater than 30", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 71}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrHeightOutOfRange.Error(), resp["message"])
	})

	t.Run("Return 400 when y is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 0, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrNegativeZeroBuilder("y").Error(), resp["message"])
	})

	t.Run("Return 400 when x is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 0, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)
		ec.Request().Header.Set("Content-Type", "application/json")

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, ErrNegativeZeroBuilder("x").Error(), resp["message"])
	})

	t.Run("Return 400 when header missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate", strings.NewReader("{\"x\": 2, \"y\": 1, \"height\": 11}"))
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		err := server.PostEstateIdTree(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resRecorder.Code)
		assert.Equal(t, "code=415, message=Unsupported Media Type", resp["message"])
	})
}

func TestGetEstateIdStats(t *testing.T) {
	t.Run("Return 200 when median exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			Count:  5,
			Max:    11,
			Min:    1,
			Median: 3,
		}

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, nil)

		err := server.GetEstateIdStats(ec, id)

		resp := readJson[generated.EstateStatResponse](t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resRecorder.Code)
		assert.Equal(t, estRep.Count, resp.Count)
		assert.Equal(t, estRep.Max, resp.Max)
		assert.Equal(t, estRep.Median, resp.Median)
		assert.Equal(t, estRep.Min, resp.Min)
	})

	t.Run("Return 200 when median is not exist and data is odd", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			Count:  5,
			Max:    11,
			Min:    1,
			Median: 0,
		}
		median := 2.

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, nil)
		mockRepo.EXPECT().GetHeightEstateTrees(ec.Request().Context(), repository.GetHeightEstateTreesInput{
			EstateId: id,
		}).Return(repository.GetHeightEstateTreesOutput{
			Heights: []int{3, 2, 1, 11, 2},
		}, nil)
		mockRepo.EXPECT().StoreMedianEstate(ec.Request().Context(), repository.StoreMedianEstateInput{
			EstateId: id,
			Median:   float64(median),
		}).Return(nil)

		err := server.GetEstateIdStats(ec, id)

		resp := readJson[generated.EstateStatResponse](t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resRecorder.Code)
		assert.Equal(t, estRep.Count, resp.Count)
		assert.Equal(t, estRep.Max, resp.Max)
		assert.Equal(t, median, resp.Median)
		assert.Equal(t, estRep.Min, resp.Min)
	})

	t.Run("Return 200 when median is not exist and data is even", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			Count:  8,
			Max:    11,
			Min:    1,
			Median: 0,
		}
		median := 5.5

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, nil)
		mockRepo.EXPECT().GetHeightEstateTrees(ec.Request().Context(), repository.GetHeightEstateTreesInput{
			EstateId: id,
		}).Return(repository.GetHeightEstateTreesOutput{
			Heights: []int{3, 2, 1, 11, 2, 4, 7, 9},
		}, nil)
		mockRepo.EXPECT().StoreMedianEstate(ec.Request().Context(), repository.StoreMedianEstateInput{
			EstateId: id,
			Median:   float64(median),
		}).Return(nil)

		err := server.GetEstateIdStats(ec, id)

		resp := readJson[generated.EstateStatResponse](t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resRecorder.Code)
		assert.Equal(t, estRep.Count, resp.Count)
		assert.Equal(t, estRep.Max, resp.Max)
		assert.Equal(t, median, resp.Median)
		assert.Equal(t, estRep.Min, resp.Min)
	})

	t.Run("Return 500 when median is not exist and get height trees error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			Count:  8,
			Max:    11,
			Min:    1,
			Median: 0,
		}

		errAny := errors.New("any error")

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, nil)
		mockRepo.EXPECT().GetHeightEstateTrees(ec.Request().Context(), repository.GetHeightEstateTreesInput{
			EstateId: id,
		}).Return(repository.GetHeightEstateTreesOutput{
			Heights: []int{3, 2, 1, 11, 2, 4, 7, 9},
		}, errAny)

		err := server.GetEstateIdStats(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resRecorder.Code)
		assert.Equal(t, errAny.Error(), resp["message"])
	})

	t.Run("Return 500 when median is not exist and get estate error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			Count:  8,
			Max:    11,
			Min:    1,
			Median: 0,
		}

		errAny := errors.New("any error")

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, errAny)

		err := server.GetEstateIdStats(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resRecorder.Code)
		assert.Equal(t, errAny.Error(), resp["message"])
	})

	t.Run("Return 404 when median is not exist and get estate missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			Count:  8,
			Max:    11,
			Min:    1,
			Median: 0,
		}

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, sql.ErrNoRows)

		err := server.GetEstateIdStats(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resRecorder.Code)
		assert.Equal(t, ErrNotFoundBuilder("estate").Error(), resp["message"])
	})
}

func TestGetEstateIdDronePlan(t *testing.T) {
	t.Run("Return 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plan", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			DroneDistance: 15,
		}

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, nil)

		err := server.GetEstateIdDronePlan(ec, id)

		resp := readJson[generated.EstateDronePlanResponse](t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resRecorder.Code)
		assert.Equal(t, estRep.DroneDistance, resp.Distance)
	})

	t.Run("Return 500 when get estate error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plan", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			Count:  8,
			Max:    11,
			Min:    1,
			Median: 0,
		}

		errAny := errors.New("any error")

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, errAny)

		err := server.GetEstateIdDronePlan(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resRecorder.Code)
		assert.Equal(t, errAny.Error(), resp["message"])
	})

	t.Run("Return 404 when median is not exist and get estate missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plane", nil)
		resRecorder := httptest.NewRecorder()

		ec := echo.New().NewContext(req, resRecorder)

		mockRepo := repository.NewMockRepositoryInterface(ctrl)

		server := Server{
			Repository: mockRepo,
		}

		id := uuid.New().String()

		estRep := repository.GetEstateByIdOutput{
			Count:  8,
			Max:    11,
			Min:    1,
			Median: 0,
		}

		mockRepo.EXPECT().GetEstateById(ec.Request().Context(), repository.GetEstateByIdInput{
			Id: id,
		}).Return(estRep, sql.ErrNoRows)

		err := server.GetEstateIdDronePlan(ec, id)

		resp := readJsonResult(t, resRecorder.Result())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resRecorder.Code)
		assert.Equal(t, ErrNotFoundBuilder("estate").Error(), resp["message"])
	})
}

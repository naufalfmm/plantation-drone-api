package handler

import (
	"database/sql"
	"math"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/naufalfmm/plantation-drone-api/generated"
	"github.com/naufalfmm/plantation-drone-api/repository"
)

// (POST /estate)
func (s *Server) PostEstate(ctx echo.Context) error {
	var req generated.CreateEstateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if req.Length <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: ErrNegativeZeroBuilder("length").Error(),
		})
	}

	if req.Width <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: ErrNegativeZeroBuilder("width").Error(),
		})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	id := uuid.New().String()
	err := s.Repository.CreateEstate(ctx.Request().Context(), repository.CreateEstateInput{
		Id:     id,
		Width:  req.Width,
		Length: req.Length,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, generated.UuidResponse{
		Id: id,
	})
}

// The endpoint of storing tree in specific point of the estate
// (POST /estate/{id}/tree)
func (s *Server) PostEstateIdTree(ctx echo.Context, id string) error {
	var req generated.CreateTreeRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if req.X <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: ErrNegativeZeroBuilder("x").Error(),
		})
	}

	if req.Y <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: ErrNegativeZeroBuilder("y").Error(),
		})
	}

	if req.Height <= 0 || req.Height > 30 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: ErrHeightOutOfRange.Error(),
		})
	}

	est, err := s.Repository.GetEstateById(ctx.Request().Context(), repository.GetEstateByIdInput{
		Id: id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
				Message: ErrNotFoundBuilder("estate").Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if req.X > est.Length || req.Y > est.Width {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: ErrCoordinateOutOfBound.Error(),
		})
	}

	c, err := s.Repository.CountCoordinateTree(ctx.Request().Context(), repository.CountCoordinateTreeInput{
		X: req.X,
		Y: req.Y,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}
	if c.Count > 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: ErrTreeExist.Error(),
		})
	}

	prevX, prevY, nextX, nextY := getPrevNextCoordinate(req.X, req.Y, est.Length)

	prevNextHeights, err := s.Repository.GetPrevNextTree(ctx.Request().Context(), repository.GetPrevNextTreeInput{
		PrevX: prevX,
		PrevY: prevY,
		NextX: nextX,
		NextY: nextY,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	treeId := uuid.New().String()
	err = s.Repository.CreateTree(ctx.Request().Context(), repository.CreateTreeInput{
		Id:     treeId,
		X:      req.X,
		Y:      req.Y,
		Height: req.Height,

		EstateId:        id,
		DroneDistFactor: int(math.Abs(float64(req.Height)-float64(prevNextHeights.PrevTreeHeight)) + math.Abs(float64(req.Height)-float64(prevNextHeights.NextTreeHeight)) - float64(prevNextHeights.PrevTreeHeight) - float64(prevNextHeights.NextTreeHeight)),
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, generated.UuidResponse{
		Id: treeId,
	})
}

// The endpoint of retrieving the estate stats, that are max, min, count, and median of trees
// (GET /estate/{id}/stats)
func (s *Server) GetEstateIdStats(ctx echo.Context, id string) error {
	est, err := s.Repository.GetEstateById(ctx.Request().Context(), repository.GetEstateByIdInput{
		Id: id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
				Message: ErrNotFoundBuilder("estate").Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	median := est.Median
	if median == 0 {

		treeHeights, err := s.Repository.GetHeightEstateTrees(ctx.Request().Context(), repository.GetHeightEstateTreesInput{
			EstateId: id,
		})
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
				Message: err.Error(),
			})
		}

		median = findMedian(treeHeights.Heights)

		go s.Repository.StoreMedianEstate(ctx.Request().Context(), repository.StoreMedianEstateInput{
			EstateId: id,
			Median:   median,
		})
	}

	return ctx.JSON(http.StatusCreated, generated.EstateStatResponse{
		Count:  est.Count,
		Max:    est.Max,
		Median: median,
		Min:    est.Min,
	})
}

// The endpoint of retrieving the estate drone plan
// (GET /estate/{id}/drone-plan)
func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id string, params generated.GetEstateIdDronePlanParams) error {
	est, err := s.Repository.GetEstateById(ctx.Request().Context(), repository.GetEstateByIdInput{
		Id: id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
				Message: ErrNotFoundBuilder("estate").Error(),
			})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, generated.EstateDronePlanResponse{
		Distance: est.DroneDistance,
	})
}

// func (s *Server) getRestEstateIdDronePlan(ctx echo.Context, id string, length, maxDistance int) (x int, y int, err error) {
// 	estTrees, err := s.Repository.GetAllEstateTrees(ctx.Request().Context(), repository.GetAllEstateTreesInput{
// 		EstateId: id,
// 		Orders:   []string{"y", "x"},
// 		Limit:    int((maxDistance-2)/10) + 1,
// 	})
// 	if err != nil {
// 		return
// 	}

// 	for _, estTree := range estTrees.EstateTrees {
// 	}

// 	return
// }

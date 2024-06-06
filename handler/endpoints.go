package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/naufalfmm/plantation-drone-api/generated"
	"github.com/naufalfmm/plantation-drone-api/repository"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) GetHello(ctx echo.Context, params generated.GetHelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

// (POST /estate)
func (s *Server) PostEstate(ctx echo.Context) error {
	var req generated.CreateEstateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
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

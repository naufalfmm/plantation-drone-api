package main

import (
	"os"

	"github.com/naufalfmm/plantation-drone-api/generated"
	"github.com/naufalfmm/plantation-drone-api/handler"
	"github.com/naufalfmm/plantation-drone-api/helper"
	"github.com/naufalfmm/plantation-drone-api/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	e.Validator = helper.NewValidator()

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}

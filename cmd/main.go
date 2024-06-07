package main

import (
	"os"

	"github.com/naufalfmm/plantation-drone-api/generated"
	"github.com/naufalfmm/plantation-drone-api/handler"
	"github.com/naufalfmm/plantation-drone-api/repository"
	"github.com/naufalfmm/plantation-drone-api/utils/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	database, err := db.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Db: database,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}

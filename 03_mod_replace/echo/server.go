package main

import (
	"github.com/equilibristofgo/sandbox/03_mod_replace/pkg"

	"github.com/labstack/echo/v4"
)

var (
	pageSize  = 4
	apiPrefix = "/api/"
)

type WebApplication struct {
	echoApp          *echo.Echo
	entityRepository pkg.EntityRepository
}

func NewWebApplication() *WebApplication {
	er := pkg.NewEntityMemoryRepository()
	er.Init()

	return &WebApplication{
		echoApp:          echo.New(),
		entityRepository: er,
	}
}

func (wa *WebApplication) Init() {
}

func (wa *WebApplication) Start(address string) {
	wa.echoApp.Logger.Fatal(wa.echoApp.Start(address))
}

func main() {
	wa := NewWebApplication()

	wa.Init()
	wa.Start(":8080")
}

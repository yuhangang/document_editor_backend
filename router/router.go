package router

import (
	"echoapp/container"
	"echoapp/controller"
	"echoapp/repo"

	_ "echoapp/docs" // for using echo-swagger

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, container container.Container, repo *repo.Repo) {
	setLocationController(e, container, repo)
	// setSwagger(container, e)
}

func setLocationController(e *echo.Echo, container container.Container, repo *repo.Repo) {
	ph := controller.NewContinentHandler(e.StdLogger, container)
	e.GET("/continents",
		func(c echo.Context) error { return ph.GetContinents(c) })
	//e.GET("/countries", ph.GetCountries)
}

//func setSwagger(container container.Container, e *echo.Echo) {
//	if container.GetEnv() == config.DEV {
//		e.GET("/swagger/*", echoSwagger.WrapHandler)
//	}
// }

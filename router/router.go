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
	locationController := controller.NewLocationController(e.StdLogger, container)
	userController := controller.NewUserController(e.StdLogger, container)
	documentFileController := controller.NewDocumentFileController(e.StdLogger, container)

	e.GET("/continents",
		func(c echo.Context) error { return locationController.GetContinents(c) })
	e.GET("/countries", func(c echo.Context) error { return locationController.GetCountries(c) })
	e.GET("/cities", func(c echo.Context) error { return locationController.GetCities(c) })

	e.GET("/devices", func(c echo.Context) error { return userController.GetDevices(c) })
	e.POST("/devices", func(c echo.Context) error { return userController.CreateUserDevice(c) })

	e.GET("/documents", func(c echo.Context) error { return documentFileController.GetDocumentFiles(c) })
	e.POST("/documents", func(c echo.Context) error { return documentFileController.CreateDocumentFile(c) })
	e.PUT("/documents", func(c echo.Context) error { return documentFileController.UpdateDocumentFile(c) })

	// initial invoke cache for master data
	locationController.LoadMasterData()
}

//func setSwagger(container container.Container, e *echo.Echo) {
//	if container.GetEnv() == config.DEV {
//		e.GET("/swagger/*", echoSwagger.WrapHandler)
//	}
// }

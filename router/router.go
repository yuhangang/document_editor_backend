package router

import (
	"echoapp/container"
	"echoapp/controller"
	"echoapp/repo"

	_ "echoapp/docs" // for using echo-swagger

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo, container container.Container, repo *repo.Repo) {
	setLocationController(e, container, repo)
	// setSwagger(container, e)
}

func setLocationController(e *echo.Echo, container container.Container, repo *repo.Repo) {
	locationController := controller.NewLocationController(e.StdLogger, container)
	userController := controller.NewUserController(e.StdLogger, container)
	documentFileController := controller.NewDocumentFileController(e.StdLogger, container)
	authController := controller.NewAuthController(e.StdLogger, container)

	e.GET("/continents", locationController.GetContinents)
	e.GET("/countries", locationController.GetCountries)
	e.GET("/cities", locationController.GetCities)

	e.GET("/devices", userController.GetDevices)
	e.POST("/devices", userController.CreateUserDevice)

	e.GET("/documents", documentFileController.GetDocumentFiles)
	e.POST("/documents", documentFileController.CreateDocumentFile)
	e.PUT("/documents", documentFileController.UpdateDocumentFile)

	e.POST("/auth", authController.Login)

	e.GET("/", authController.Accessible)
	r := e.Group("/api")
	// Configure middleware with the custom claims type
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	jwtConfig := middleware.JWTConfig{
		Claims:     &controller.JwtDevice{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(jwtConfig))
	r.GET("", authController.Restricted)
	r.GET("/documents", documentFileController.GetDocumentFiles)
	r.POST("/documents", documentFileController.CreateDocumentFile)
	r.PUT("/documents", documentFileController.UpdateDocumentFile)

	// initial invoke cache for master data
	locationController.LoadMasterData()
}

//func setSwagger(container container.Container, e *echo.Echo) {
//	if container.GetEnv() == config.DEV {
//		e.GET("/swagger/*", echoSwagger.WrapHandler)
//	}
// }

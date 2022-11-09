package router

import (
	"echoapp/container"
	"echoapp/controller"
	_ "echoapp/docs" // for using echo-swagger
	"echoapp/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo, container container.Container) {
	setControllers(e, &container)
	// setSwagger(container, e)
}

func setControllers(e *echo.Echo, container *container.Container) {
	locationController := controller.NewLocationController(e.StdLogger, *container)
	userController := controller.NewUserController(e.StdLogger, *container)
	documentFileController := controller.NewDocumentFileController(e.StdLogger, *container)
	authController := controller.NewAuthController(e.StdLogger, *container)

	middleware.InitLoggerMiddleware(e, *container)

	e.GET("/continents", locationController.GetContinents)
	e.GET("/countries", locationController.GetCountries)
	e.GET("/cities", locationController.GetCities)

	e.GET("/devices", userController.GetDevices)
	e.POST("/devices", userController.CreateUserDevice)

	e.POST("/auth", authController.Login)

	e.GET("/", authController.Accessible)
	r := e.Group("/api")
	// Configure middleware with the custom claims type
	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	jwtConfig := echoMiddleware.JWTConfig{
		Claims:     &controller.JwtDevice{},
		SigningKey: []byte("secret"),
	}
	r.Use(echoMiddleware.JWTWithConfig(jwtConfig))
	r.GET("", authController.Restricted)
	r.GET("/documents", documentFileController.GetDocumentFiles)
	r.POST("/documents", documentFileController.CreateDocumentFile)
	r.GET("/documents/:id", documentFileController.GetDocumentFileById)
	r.PUT("/documents/:id", documentFileController.UpdateDocumentFile)
	r.DELETE("/documents/:id", documentFileController.DeleteDocumentFile)

	// initial invoke cache for master data
	locationController.LoadMasterData()
}

//func setSwagger(container container.Container, e *echo.Echo) {
//	if container.GetEnv() == config.DEV {
//		e.GET("/swagger/*", echoSwagger.WrapHandler)
//	}
// }

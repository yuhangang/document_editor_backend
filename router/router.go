package router

import (
	"echoapp/container"
	"echoapp/controller"
	_ "echoapp/docs" // for using echo-swagger
	"echoapp/middleware"
	"os"

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
	fileController := controller.NewFileController(e.StdLogger, *container)

	middleware.InitLoggerMiddleware(e, *container)

	e.GET("/continents", locationController.GetContinents)
	e.GET("/countries", locationController.GetCountries)
	e.GET("/cities", locationController.GetCities)

	e.GET("/devices", userController.GetDevices)
	e.POST("/devices", userController.CreateUserDevice)

	e.POST("/auth", authController.Login)

	e.GET("/", authController.Accessible)
	api := e.Group("/api")
	// Configure middleware with the custom claims type
	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	jwtConfig := echoMiddleware.JWTConfig{
		Claims:     &controller.JwtDevice{},
		SigningKey: []byte("secret"),
	}
	api.Use(echoMiddleware.JWTWithConfig(jwtConfig))
	api.GET("", authController.Restricted)
	api.GET("/documents", documentFileController.GetDocumentFiles)
	api.POST("/documents", documentFileController.CreateDocumentFile)
	api.GET("/documents/:id", documentFileController.GetDocumentFileById)
	api.PUT("/documents/:id", documentFileController.UpdateDocumentFile)
	api.DELETE("/documents/:id", documentFileController.DeleteDocumentFile)

	if os.MkdirAll(controller.FileDirectory, 0777) != nil {
		panic("Unable to create directory for tagfile!")
	}
	e.GET("/attachments/:id", fileController.GetAttachment)
	api.Static("/attachments/upload", "public/upload")
	api.POST("/attachments", fileController.UploadAttachment)
	api.DELETE("/attachments/:id", fileController.DeleteAttachment)

	// initial invoke cache for master data
	locationController.LoadMasterData()
}

//func setSwagger(container container.Container, e *echo.Echo) {
//	if container.GetEnv() == config.DEV {
//		e.GET("/swagger/*", echoSwagger.WrapHandler)
//	}
// }

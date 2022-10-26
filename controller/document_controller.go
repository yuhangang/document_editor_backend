package controller

import (
	"echoapp/container"
	"echoapp/model"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DocumentFileController interface {
	CreateDocumentFile(c echo.Context) error
	GetDocumentFiles(c echo.Context) error
	UpdateDocumentFile(c echo.Context) error
}

// Products is a http.Handler
type documentFileController struct {
	l         *log.Logger
	container container.Container
}

func NewDocumentFileController(l *log.Logger, container container.Container) DocumentFileController {
	return &documentFileController{l, container}
}

func (u *documentFileController) CreateDocumentFile(c echo.Context) error {
	callback := c.QueryParam("callback")
	var newDocument model.DocumentFile

	err := c.Bind(&newDocument)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}

	insertError := u.container.GetRepo().DB.Create(&newDocument).Error
	if insertError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &newDocument)
}

func (u *documentFileController) GetDocumentFiles(c echo.Context) error {
	callback := c.QueryParam("callback")
	var documentFiles []model.DocumentFile
	err := u.container.GetRepo().DB.Find(&documentFiles).Error

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &documentFiles)
}

func (u *documentFileController) UpdateDocumentFile(c echo.Context) error {
	callback := c.QueryParam("callback")
	var newDocument model.DocumentFile

	err := c.Bind(&newDocument)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}

	insertError := u.container.GetRepo().DB.Create(&newDocument).Error
	if insertError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &newDocument)
}

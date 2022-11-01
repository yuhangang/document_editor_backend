package controller

import (
	"echoapp/constant"
	"echoapp/container"
	"echoapp/model"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

func (u *documentFileController) getDB() *gorm.DB {
	return u.container.GetRepo().DB
}

func NewDocumentFileController(l *log.Logger, container container.Container) DocumentFileController {
	return &documentFileController{l, container}
}

func (u *documentFileController) decryptJWT(c echo.Context) (JwtDevice, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtDevice)
	return *claims, nil
}

func (u *documentFileController) CreateDocumentFile(c echo.Context) error {
	var newDocument model.DocumentFile
	device, jwtErr := u.decryptJWT(c)
	if jwtErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	err := c.Bind(&newDocument)
	if err != nil {
		u.container.GetLogger().GetZapLogger().Debugln(err)
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}
	newDocument.DeviceId = device.DeviceId

	insertError := u.getDB().Create(&newDocument).Error
	if insertError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}
	return c.JSON(http.StatusOK, &newDocument)
}

func (u *documentFileController) GetDocumentFiles(c echo.Context) error {
	// callback := c.QueryParam("callback")
	device, jwtErr := u.decryptJWT(c)
	if jwtErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	var documentFiles []model.DocumentFile
	err := u.getDB().Where("device_id = ?", device.DeviceId).Find(&documentFiles).Error

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}
	res := DocumentListResponse{Data: documentFiles}

	resJson, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}
	return c.JSONBlob(http.StatusOK, resJson)
}

type DocumentListResponse struct {
	Data []model.DocumentFile `json:"data"`
}

func (u *documentFileController) UpdateDocumentFile(c echo.Context) error {
	docId := c.QueryParam("id")
	if len(docId) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, constant.InternalServerErrorMsg)
	}
	var columnsToUpdate interface{}

	var exists bool
	checkExistsError := u.getDB().Model(&model.DocumentFile{}).
		Select("count(*) > 0").
		Where("id = ?", docId).
		Find(&exists).
		Error

	if checkExistsError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Document not found")
	}

	err := c.Bind(&columnsToUpdate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	var document model.DocumentFile

	insertError := u.getDB().Model(&document).Where("id = ?", docId).UpdateColumns(&columnsToUpdate).Error
	if insertError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	return c.JSON(http.StatusOK, &document)
}

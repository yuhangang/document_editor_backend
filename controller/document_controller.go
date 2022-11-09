package controller

import (
	constant "echoapp/commons"
	"echoapp/container"
	"echoapp/model"
	"echoapp/repo"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type DocumentFileController interface {
	CreateDocumentFile(c echo.Context) error
	GetDocumentFiles(c echo.Context) error
	UpdateDocumentFile(c echo.Context) error
	DeleteDocumentFile(c echo.Context) error
	GetDocumentFileById(c echo.Context) error
}

// Products is a http.Handler
type documentFileController struct {
	l         *log.Logger
	container container.Container
}

func (u *documentFileController) getDB() repo.Repository {
	return u.container.GetRepository()
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

	var uuidErr error
	documentId := uuid.Must(uuid.NewV4(), uuidErr).String()
	if uuidErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}
	newDocument.DeviceId = device.DeviceId
	newDocument.DocumentId = documentId

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
	docId := c.Param("id")
	if len(docId) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, constant.InternalServerErrorMsg)
	}
	var columnsToUpdate interface{}

	var exists bool
	checkExistsError := u.getDB().Model(&model.DocumentFile{}).
		Select("count(*) > 0").
		Where("document_id = ?", docId).
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

	insertError := u.getDB().Model(&model.DocumentFile{}).Where("document_id = ?", docId).UpdateColumns(&columnsToUpdate).Error
	if insertError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	var documentFile model.DocumentFile
	queryError := u.getDB().
		Where("document_id = ?", docId).
		First(&documentFile).
		Error

	if queryError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	return c.JSON(http.StatusOK, &documentFile)
}

func (u *documentFileController) GetDocumentFileById(c echo.Context) error {
	docId := c.Param("id")
	if len(docId) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, constant.InternalServerErrorMsg)
	}

	var documentFile model.DocumentFile
	queryError := u.getDB().
		Where("document_id = ?", docId).
		First(&documentFile).
		Error

	if queryError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	return c.JSON(http.StatusOK, &documentFile)
}

func (u *documentFileController) DeleteDocumentFile(c echo.Context) error {
	docId := c.Param("id")
	device, jwtErr := u.decryptJWT(c)
	if jwtErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}
	if len(docId) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, constant.InternalServerErrorMsg)
	}

	var documentFile model.DocumentFile
	queryError := u.getDB().
		Where("document_id = ?", docId).
		First(&documentFile).
		Error
	if device.DeviceId != documentFile.DeviceId {
		return echo.NewHTTPError(http.StatusUnauthorized, constant.UnauthorizedErrorMsg)
	}

	if queryError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	deletionError := u.getDB().Delete(&documentFile, documentFile.ID).Error
	if deletionError != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}

	return c.JSON(http.StatusOK, "Document deleted.")
}

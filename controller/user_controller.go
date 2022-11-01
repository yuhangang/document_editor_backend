package controller

import (
	"echoapp/constant"
	"echoapp/container"
	"echoapp/model"
	"echoapp/model/dto"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController interface {
	CreateUserDevice(c echo.Context) error
	GetDevices(c echo.Context) error
}

// Products is a http.Handler
type userController struct {
	l         *log.Logger
	container container.Container
}

func NewUserController(l *log.Logger, container container.Container) UserController {
	return &userController{l, container}
}

func (u *userController) CreateUserDevice(c echo.Context) error {
	dto := dto.NewDeviceInfoDto()
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if errors := dto.Validate(); errors != nil {
		return c.JSON(http.StatusBadRequest, errors)
	}

	deviceInfo := dto.Create()
	tx := u.container.GetRepo().DB.Begin()
	var dbInsertError error
	var user model.DeviceInfo
	dbErr := tx.Take(&user).Error
	if dbErr != nil && errors.Is(dbErr, gorm.ErrRecordNotFound) {
		dbInsertError = tx.Create(&deviceInfo).Error
	} else {
		dbInsertError = tx.Model(&deviceInfo).
			Where("device_id = ?", deviceInfo.DeviceId).
			Update("device_os_version", "device_os_version").Error
	}

	if dbInsertError != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, dbInsertError)
	}
	var newDeviceInfo model.DeviceInfo
	queryError := tx.Where("device_id = ?", deviceInfo.DeviceId).First(&newDeviceInfo).Error
	if queryError != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, queryError)
	}
	tx.Commit()
	return c.JSON(http.StatusOK, newDeviceInfo)
}

func (u *userController) GetDevices(c echo.Context) error {
	callback := c.QueryParam("callback")
	var deviceInfo []model.DeviceInfo
	err := u.container.GetRepo().DB.Find(&deviceInfo).Error

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, constant.InternalServerErrorMsg)
	}
	return c.JSONP(http.StatusOK, callback, &deviceInfo)
}

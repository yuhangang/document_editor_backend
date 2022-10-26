package controller

import (
	"echoapp/container"
	"echoapp/model"
	"echoapp/model/dto"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
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
		fmt.Println("yolo 1", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	if errors := dto.Validate(); errors != nil {
		fmt.Println("yolo 2", errors)
		return c.JSON(http.StatusBadRequest, errors)
	}

	deviceInfo := dto.Create()
	dbInsertError := u.container.GetRepo().DB.Create(&deviceInfo).Error
	if dbInsertError != nil {
		fmt.Println("yolo 3", dbInsertError)
		return c.JSON(http.StatusBadRequest, dbInsertError)
	}
	return c.JSON(http.StatusOK, deviceInfo)
}

func (u *userController) GetDevices(c echo.Context) error {
	callback := c.QueryParam("callback")
	var deviceInfo []model.DeviceInfo
	err := u.container.GetRepo().DB.Find(&deviceInfo).Error

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &deviceInfo)
}

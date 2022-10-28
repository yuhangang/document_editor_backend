package controller

import (
	"echoapp/container"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AuthController interface {
	Login(c echo.Context) error
	Restricted(c echo.Context) error
	Accessible(c echo.Context) error
}

type JwtDevice struct {
	DeviceId string `json:"device_id"`
	jwt.StandardClaims
}

type authController struct {
	l         *log.Logger
	container container.Container
}

func NewAuthController(l *log.Logger, container container.Container) AuthController {
	return &authController{l, container}
}

func (a *authController) Login(c echo.Context) error {
	deviceId := c.FormValue("device_id")

	// Throws unauthorized error
	if len(deviceId) == 0 {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &JwtDevice{
		deviceId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func (a *authController) Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func (a *authController) Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtDevice)
	name := claims.DeviceId
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

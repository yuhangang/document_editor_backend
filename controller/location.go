package controller

import (
	"echoapp/container"
	"echoapp/model"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ContinentController interface {
	GetContinents(c echo.Context) error
	GetCountries(c echo.Context) error
	GetCities(c echo.Context) error
}

// Products is a http.Handler
type continentController struct {
	l         *log.Logger
	container container.Container
}

func NewContinentHandler(l *log.Logger, container container.Container) ContinentController {
	return &continentController{l, container}
}

func (p *continentController) GetContinents(c echo.Context) error {
	callback := c.QueryParam("callback")
	var continents []model.Continent
	err := p.container.GetRepo().DB.Find(&continents).Error

	if err != nil {
		p.l.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &continents)
}

func (p *continentController) GetCountries(c echo.Context) error {
	callback := c.QueryParam("callback")
	var countries []model.Country
	err := p.container.GetRepo().DB.Find(&countries).Error

	if err != nil {
		p.l.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &countries)
}

func (p *continentController) GetCities(c echo.Context) error {
	callback := c.QueryParam("callback")
	var cities []model.City
	err := p.container.GetRepo().DB.Find(&cities).Error

	if err != nil {
		p.l.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &cities)
}

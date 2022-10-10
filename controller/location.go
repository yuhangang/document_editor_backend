package controller

import (
	"echoapp/container"
	"echoapp/model"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LocationController interface {
	GetContinents(c echo.Context) error
	GetCountries(c echo.Context) error
	GetCities(c echo.Context) error
}

// Products is a http.Handler
type locationController struct {
	l         *log.Logger
	container container.Container
}

func NewLocationController(l *log.Logger, container container.Container) LocationController {
	return &locationController{l, container}
}

func (p *locationController) GetContinents(c echo.Context) error {
	callback := c.QueryParam("callback")
	var continents []model.Continent
	err := p.container.GetRepo().DB.Model(&model.Continent{}).Preload("Countries.Cities").Find(&continents).Error
	// err := p.container.GetRepo().DB.Table("Continents").
	// 	Joins("INNER JOIN Countries c ON c.continent_id = Continents.code").
	// 	Select("Continents.code, Continents.name").
	// 	Find(&continents)

	if err != nil {
		p.l.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}

	return c.JSONP(http.StatusOK, callback, &continents)
}

func (p *locationController) GetCountries(c echo.Context) error {
	callback := c.QueryParam("callback")
	var countries []model.Country
	err := p.container.GetRepo().DB.Model(&model.Country{}).Preload("Cities").Find(&countries).Error

	if err != nil {
		p.l.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &countries)
}

func (p *locationController) GetCities(c echo.Context) error {
	callback := c.QueryParam("callback")
	var cities []model.City
	err := p.container.GetRepo().DB.Find(&cities).Error

	if err != nil {
		p.l.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Something wrong with internal server please try again later")
	}
	return c.JSONP(http.StatusOK, callback, &cities)
}

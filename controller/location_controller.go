package controller

import (
	"echoapp/constant"
	"echoapp/container"
	"echoapp/model"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LocationController interface {
	GetContinents(c echo.Context) error
	GetCountries(c echo.Context) error
	GetCities(c echo.Context) error
	LoadMasterData() ([]model.Continent, error)
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

	entry, cacheError := p.container.GetBigCache().Get(constant.Cache_continent_key)
	if cacheError != nil {
		p.l.Fatal(cacheError)

	} else {
		return c.JSONBlob(http.StatusOK, entry)
	}

	continents, err := p.LoadMasterData()

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

func (p *locationController) LoadMasterData() ([]model.Continent, error) {
	// Pre cached query of master data

	var continents []model.Continent
	queryError := p.container.GetRepo().DB.Model(&model.Continent{}).Preload("Countries.Cities").Find(&continents).Error
	if queryError == nil {
		jsonByte, jsonError := json.Marshal(continents)
		if jsonError == nil {

			p.container.GetBigCache().Set(constant.Cache_continent_key, jsonByte)
			return continents, nil
		} else {

			return nil, queryError
		}
	} else {

		return nil, queryError
	}

}
package data_processing

import (
	"echoapp/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ExecuteProcessing() error {
	continents, processContinentDataError := processContinentData()
	if processContinentDataError != nil {
		return processContinentDataError
	}

	writeContinentDataError := writeDataToJsonFile(continents, "data/processed/continents.json")
	if writeContinentDataError != nil {
		return writeContinentDataError
	}

	cities, processCityDataError := processCityData()
	if processCityDataError != nil {
		return processCityDataError
	}
	writeCitiesDataError := writeDataToJsonFile(cities, "data/processed/world_cities.json")
	if writeCitiesDataError != nil {
		return writeCitiesDataError
	}

	countries, processCountryDataError := processCountryData()
	if processCountryDataError != nil {
		return processCountryDataError
	}

	writeCountryDataError := writeDataToJsonFile(countries, "data/processed/world_countries.json")
	if writeCountryDataError != nil {
		return writeCitiesDataError
	}
	return nil

}

func writeDataToJsonFile(data any, directory string) error {
	file, jsonParseError := json.MarshalIndent(data, "", " ")
	if jsonParseError != nil {
		return jsonParseError
	}
	writeJsonFileError := ioutil.WriteFile(directory, file, 0644)
	if writeJsonFileError != nil {
		return writeJsonFileError
	}
	return nil
}

func processContinentData() ([]model.Continent, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("data/continents.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	data := make(map[string]string)
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}
	var continents []model.Continent
	for key, element := range data {
		continent := model.Continent{Code: key, Name: element}
		continents = append(continents, continent)
	}

	fmt.Println("Successfully Opened continents.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	return continents, nil
}

func processCityData() ([]model.City, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("data/world_cities.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var data []map[string]interface{}
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}
	//fmt.Println(data)
	cities := []model.City{}
	for _, json := range data {
		city := model.City{
			Name:      json["city_ascii"].(string),
			CountryID: json["iso2"].(string),
			Capital:   json["capital"].(string),
			Lat:       json["lat"].(float64),
			Lng:       json["lat"].(float64)}
		cities = append(cities, city)
	}

	fmt.Println("Successfully Opened world_cities.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	return cities, nil
}

func processCountryData() ([]model.Country, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("data/countries.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	data := make(map[string]interface{})
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println("panic!")
		return nil, err
	}

	var countries []model.Country
	for code, jsonStr := range data {
		data := jsonStr.(map[string]interface{})
		country := model.Country{
			Code:        code,
			Name:        data["name"].(string),
			Native:      data["native"].(string),
			Capital:     data["capital"].(string),
			ContinentID: data["continent"].(string)}
		countries = append(countries, country)
	}

	fmt.Println("Successfully Opened countries.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	return countries, nil
}

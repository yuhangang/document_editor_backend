package migration

import (
	"echoapp/container"
	"echoapp/model"
	"encoding/json"
	"io/ioutil"
	"os"

	"gorm.io/gorm"
)

// InitMasterData creates the master data used in this application.
func InitMasterData(container container.Container) {
	errs := make(chan error)
	if container.GetConfig().Extension.MasterGenerator {

		initContinentsErr := initContinentData(container, errs)
		if initContinentsErr != nil {
			container.GetLogger().GetZapLogger().Fatal(initContinentsErr)
		}
		initCountriesErr := initCountriesData(container, errs)
		if initCountriesErr != nil {
			container.GetLogger().GetZapLogger().Fatal(initCountriesErr)
		}
		initCitiesErr := initCitiesData(container, errs)
		if initCitiesErr != nil {

			container.GetLogger().GetZapLogger().Fatal(1000, initCitiesErr)

		}

		//rep := container.GetRepository()

		// r := model.NewAuthority("Admin")
		// _, _ = r.Create(rep)
		// a := model.NewAccountWithPlainPassword("test", "test", r.ID)
		// _, _ = a.Create(rep)
		// a = model.NewAccountWithPlainPassword("test2", "test2", r.ID)
		// _, _ = a.Create(rep)
		//
		// c := model.NewCategory("Technical Book")
		// _, _ = c.Create(rep)
		// c = model.NewCategory("Magazine")
		// _, _ = c.Create(rep)
		// c = model.NewCategory("Novel")
		// _, _ = c.Create(rep)
		//
		// f := model.NewFormat("Paper Book")
		// _, _ = f.Create(rep)
		// f = model.NewFormat("e-Book")
		// _, _ = f.Create(rep)
	}
}

func initContinentData(container container.Container, errs chan error) error {
	jsonFile, err := os.Open("data/processed/continents.json")
	if err != nil {
		return err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var continents []model.Continent
	err = json.Unmarshal(byteValue, &continents)
	container.GetRepository().Create(&continents)
	return nil
}

func initCountriesData(container container.Container, errs chan error) error {
	jsonFile, err := os.Open("data/processed/countries.json")
	if err != nil {
		return err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var countries []model.Country
	err = json.Unmarshal(byteValue, &countries)
	container.GetRepository().Create(&countries)
	return nil
}

func initCitiesData(container container.Container, errs chan error) error {
	jsonFile, err := os.Open("data/processed/cities.json")
	if err != nil {
		return err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var cities []model.City
	err = json.Unmarshal(byteValue, &cities)

	const size = 100
	var j int
	for i := 0; i < len(cities); i += size {
		j += size
		if j > len(cities) {
			j = len(cities)
		}
		db := container.GetRepository().Session(&gorm.Session{CreateBatchSize: len(cities[i:j])})
		db.Create(cities[i:j])
	}

	return nil
}

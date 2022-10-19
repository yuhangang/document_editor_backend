package model

type Country struct {
	Code        string `gorm:"primaryKey" json:"code"`
	Name        string `gorm:"unique" json:"name"`
	Capital     string `json:"capital"`
	Native      string `json:"native"`
	ContinentID string `gorm:"foreignkey:continent_id" json:"continentID"`
	//Continent   Continent
	Cities []City `json:"cities"`
}

type City struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
	//Country   Country
	CountryID string  `gorm:"foreignkey:country_id" json:"countryID"`
	Capital   string  `json:"capital"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

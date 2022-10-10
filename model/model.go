package model

type Country struct {
	Code        string `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Capital     string `gorm:"unique"`
	Native      string
	ContinentID string `gorm:"foreignkey:continent_id"`
	//Continent   Continent
	Cities []City `json:"cities"`
}

type City struct {
	Name string `gorm:"primaryKey"`
	//Country   Country
	CountryID string `gorm:"foreignkey:country_id"`
	Capital   string
	Lat       float64
	Lng       float64
}

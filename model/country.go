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

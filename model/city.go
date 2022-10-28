package model

type City struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
	//Country   Country
	CountryID string  `gorm:"foreignkey:country_id" json:"countryID"`
	Capital   string  `json:"capital"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

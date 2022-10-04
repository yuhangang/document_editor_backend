package model

type Country struct {
	Code        string `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Capital     string `gorm:"unique"`
	Native      string
	ContinentID string
	Continent   Continent `gorm:"references:Code"`
}

type City struct {
	Name      string `gorm:"primaryKey"`
	CountryID string
	Country   Country `gorm:"references:Code"`
	lat       float32
	lng       float64
}

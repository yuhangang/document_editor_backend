package dto

import (
	"echoapp/model"
	"encoding/json"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

const (
	required string = "required"
	max      string = "max"
	min      string = "min"
)

const (
	ValidationErrMessageBookTitle string = "Please enter the title with 3 to 50 characters."
	ValidationErrMessageBookISBN  string = "Please enter the ISBN with 10 to 20 characters."
)

// DeviceInfoDto defines a data transfer object for book.
type DeviceInfoDto struct {
	DeviceId              string   `validate:"required" json:"device_id"`
	DeviceModel           string   `validate:"required" json:"device_model"`
	DeviceManufacturer    string   `validate:"required" json:"device_manufacturer"`
	DeviceOsVersion       string   `validate:"required" json:"device_os_version"`
	DeviceOsVersionNumber string   `validate:"required" json:"device_os_version_number"`
	Lat                   *float64 `json:"lat"`
	Lng                   *float64 `json:"lng"`
	DeviceUserId          *uint    `json:"device_user_id"`
}

// NewDeviceInfoDto is constructor.
func NewDeviceInfoDto() *DeviceInfoDto {
	return &DeviceInfoDto{}
}

// Create creates a book model from this DTO.
func (b *DeviceInfoDto) Create() *model.DeviceInfo {
	return &model.DeviceInfo{DeviceId: b.DeviceId,
		DeviceModel:           b.DeviceModel,
		DeviceManufacturer:    b.DeviceManufacturer,
		DeviceOsVersion:       b.DeviceOsVersion,
		DeviceOsVersionNumber: b.DeviceOsVersionNumber,
		Lat:                   b.Lat, Lng: b.Lng,
		DeviceUserId: b.DeviceUserId}
}

// Validate performs validation check for the each item.
func (b *DeviceInfoDto) Validate() map[string]string {
	return validateDto(b)
}

func validateDto(b interface{}) map[string]string {
	err := validator.New().Struct(b)
	if err == nil {
		return nil
	}

	errors := err.(validator.ValidationErrors)
	if len(errors) == 0 {
		return nil
	}

	return createErrorMessages(errors)
}

func createErrorMessages(errors validator.ValidationErrors) map[string]string {
	result := make(map[string]string)
	for i := range errors {
		switch errors[i].StructField() {
		case "Title":
			switch errors[i].Tag() {
			case required, min, max:
				result["title"] = ValidationErrMessageBookTitle
			}
		case "Isbn":
			switch errors[i].Tag() {
			case required, min, max:
				result["isbn"] = ValidationErrMessageBookISBN
			}
		default:

			result[fmt.Sprint("errors", i)] = fmt.Sprint(errors[i].Tag(), " : ", errors[i].Type())
		}
	}
	return result
}

// ToString is return string of object
func (b *DeviceInfoDto) ToString() (string, error) {
	bytes, err := json.Marshal(b)
	return string(bytes), err
}

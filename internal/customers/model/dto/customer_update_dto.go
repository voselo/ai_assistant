package dto

import (
	"time"

	"github.com/google/uuid"
)

type CustomerUpdateDTO struct {
	Id                 uuid.UUID  `json:"id" binding:"required"`
	Name               string     `json:"name"`
	Email              string     `json:"email"`
	WazzupUri          string     `json:"wazzup_uri"`
	LicenseLvl         string     `json:"license_lvl" binding:"omitempty,oneof=basic standart pro enterprise"`
	LicenseStatus      string     `json:"license_status" binding:"omitempty,oneof=active inactive freezed deleted"`
	LicenseExpiresDate *time.Time `json:"license_expires_date"`
}

func (dto *CustomerUpdateDTO) HasName() bool {
	return dto.Name != ""
}

func (dto *CustomerUpdateDTO) HasEmail() bool {
	return dto.Email != ""
}

func (dto *CustomerUpdateDTO) HasWazzupUri() bool {
	return dto.WazzupUri != ""
}

func (dto *CustomerUpdateDTO) HasLicenseLvl() bool {
	return dto.LicenseLvl != ""
}

func (dto *CustomerUpdateDTO) HasLicenseStatus() bool {
	return dto.LicenseStatus != ""
}

func (dto *CustomerUpdateDTO) HasLicenseExpiresDate() bool {
	return dto.LicenseExpiresDate != nil
}

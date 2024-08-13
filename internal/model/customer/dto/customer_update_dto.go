package customer

import (
	"time"
)

type UpdateDTO struct {
	Name               string     `json:"name"`
	Email              string     `json:"email"`
	WazzupUri          string     `json:"wazzup_uri"`
	LicenseLvl         string     `json:"license_lvl" binding:"omitempty,oneof=basic standart pro enterprise"`
	LicenseStatus      string     `json:"license_status" binding:"omitempty,oneof=active inactive freezed deleted"`
	LicenseExpiresDate *time.Time `json:"license_expires_date"`
}

func (dto *UpdateDTO) HasName() bool {
	return dto.Name != ""
}

func (dto *UpdateDTO) HasEmail() bool {
	return dto.Email != ""
}

func (dto *UpdateDTO) HasWazzupUri() bool {
	return dto.WazzupUri != ""
}

func (dto *UpdateDTO) HasLicenseLvl() bool {
	return dto.LicenseLvl != ""
}

func (dto *UpdateDTO) HasLicenseStatus() bool {
	return dto.LicenseStatus != ""
}

func (dto *UpdateDTO) HasLicenseExpiresDate() bool {
	return dto.LicenseExpiresDate != nil
}

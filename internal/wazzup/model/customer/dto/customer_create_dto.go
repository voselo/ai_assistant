package dto

import (
	model "ai_assistant/internal/model/customer"
	"time"
)

type CustomerCreateDTO struct {
	Name               string     `json:"name" binding:"required"`
	Email              string     `json:"email" binding:"required,email"`
	WazzupUri          string     `json:"wazzup_uri" binding:"required"`
	LicenseLvl         string     `json:"license_lvl" binding:"required,oneof=basic standart pro enterprise"`
	LicenseExpiresDate *time.Time `json:"license_expires_date"`
}

func (dto *CustomerCreateDTO) ToModel() *model.CustomerModel {
	return &model.CustomerModel{
		Name:               dto.Name,
		Email:              dto.Email,
		WazzupUri:          dto.WazzupUri,
		LicenseLvl:         dto.LicenseLvl,
		LicenseExpiresDate: dto.LicenseExpiresDate,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type CustomerModel struct {
	// Personal data
	Id    uuid.UUID `db:"id" json:"id"`
	Name  string    `db:"name" json:"name"`
	Email string    `db:"email" json:"email"`

	// License
	LicenseStatus      string     `db:"license_status" json:"license_status"`
	LicenseLvl         string     `db:"license_lvl" json:"license_lvl"`
	LicenseHash        string     `db:"license_hash" json:"license_hash"`
	LicenseExpiresDate *time.Time `db:"license_expires_date" json:"license_expires_date"`

	// Crm
	// CrmRefreshToken string
	// CrmAccessToken  string

	// Wazzup
	WazzupUri string `db:"wazzup_uri" json:"wazzup_uri"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

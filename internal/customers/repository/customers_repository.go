package repository

import (
	"ai_assistant/internal/customers/model"
	"ai_assistant/internal/customers/model/dto"
	"ai_assistant/pkg/logging"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type CustomersRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CustomersRepository {
	return &CustomersRepository{db: db}
}

func (repository *CustomersRepository) Create(newModel *model.CustomerModel) (*model.CustomerModel, error) {
	logger := logging.GetLogger("Info")

	newModel.CreatedAt = time.Now()
	newModel.UpdatedAt = time.Now()
	newModel.LicenseHash = generateLicenseHash(newModel.Email, newModel.LicenseLvl)
	newModel.LicenseStatus = model.LicenseStatusActive

	if newModel.LicenseExpiresDate != nil && newModel.LicenseExpiresDate.Before(time.Now()) {
		return nil, errors.New("expires date before this moment")
	}

	query := `
		INSERT INTO customers 
		(name, email, license_status, license_lvl, license_hash, license_expires_date, wazzup_uri, created_at, updated_at)
		VALUES (:name, :email, :license_status, :license_lvl, :license_hash, :license_expires_date, :wazzup_uri, :created_at, :updated_at)
		RETURNING id;
	`
	logger.Info("HASH: ", newModel.LicenseHash)
	logger.Info("Executing query: ", query, "With model: ", newModel)

	rows, err := repository.db.NamedQuery(query, newModel)
	if err != nil {
		logger.Error("Failed to insert new customer: ", err)
		return nil, fmt.Errorf("error while inserting new customer: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&newModel.Id)
		if err != nil {
			logger.Error("Failed to retrieve new customer ID: ", err)
			return nil, fmt.Errorf("error while retrieving new customer ID: %w", err)
		}
	} else {
		// Handling the case where no rows were inserted because of conflict
		logger.Info("No new customer was added, possible duplicate")
		return nil, fmt.Errorf("customer with this email or license hash already exists")
	}

	logger.Info("Customer created successfully: ", newModel.Id)
	return newModel, nil
}

func generateLicenseHash(email, licenseLevel string) string {
	// Используем текущее время для добавления уникальности
	data := fmt.Sprintf("%s:%s:%v", email, licenseLevel, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (repository *CustomersRepository) GetAll() ([]model.CustomerModel, error) {
	logger := logging.GetLogger("Info")
	var customers []model.CustomerModel
	query := "SELECT * FROM customers"
	err := repository.db.Select(&customers, query)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return customers, nil
}

func (repository *CustomersRepository) GetById(id string) (*model.CustomerModel, error) {
	query := "SELECT * FROM customers WHERE id = $1"
	var customer model.CustomerModel
	if err := repository.db.Get(&customer, query, id); err != nil {
		return nil, err
	}

	return &customer, nil

}

func (repository *CustomersRepository) Update(updateDTO *dto.CustomerUpdateDTO) (*model.CustomerModel, error) {

	var fields []string
	params := map[string]interface{}{
		"id": updateDTO.Id,
	}

	// Проверяем какие поля переданы и добавляем их в запрос
	if updateDTO.HasName() {
		fields = append(fields, "name = :name")
		params["name"] = updateDTO.Name
	}
	if updateDTO.HasEmail() {
		fields = append(fields, "email = :email")
		params["email"] = updateDTO.Email
	}
	if updateDTO.HasWazzupUri() {
		fields = append(fields, "wazzup_uri = :wazzup_uri")
		params["wazzup_uri"] = updateDTO.WazzupUri
	}
	if updateDTO.HasLicenseLvl() {
		fields = append(fields, "license_lvl = :license_lvl")
		params["license_lvl"] = updateDTO.LicenseLvl
	}
	if updateDTO.HasLicenseStatus() {
		fields = append(fields, "license_status = :license_status")
		params["license_status"] = updateDTO.LicenseStatus
	}
	if updateDTO.HasLicenseExpiresDate() {
		if updateDTO.LicenseExpiresDate.Before(time.Now()) {
			return nil, errors.New("expires date before this moment")
		}

		fields = append(fields, "license_expires_date = :license_expires_date")
		params["license_expires_date"] = updateDTO.LicenseExpiresDate
	}

	// Обновляем только если есть хотя бы одно поле для обновления
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Добавляем поле обновления времени
	fields = append(fields, "updated_at = :updated_at")
	params["updated_at"] = time.Now()

	// Строим запрос
	query := fmt.Sprintf("UPDATE customers SET %s WHERE id = :id RETURNING *", strings.Join(fields, ", "))

	// Выполняем запрос
	var updatedCustomer model.CustomerModel
	stmt, err := repository.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("could not prepare statement: %w", err)
	}
	defer stmt.Close()

	if err := stmt.Get(&updatedCustomer, params); err != nil {
		return nil, fmt.Errorf("could not execute update: %w", err)
	}

	return &updatedCustomer, nil
}

func (repository *CustomersRepository) Delete(id string) error {
	query := "DELETE FROM customers WHERE id = $1"
	_, err := repository.db.Exec(query, id)
	return err
}

func (repository *CustomersRepository) ValidateLicense(hash string) (isValid bool, err error) {
	query := "SELECT * FROM customers WHERE license_hash = $1"
	var customer model.CustomerModel
	if err := repository.db.Get(&customer, query, hash); err != nil {
		return false, err
	}

	if customer.LicenseStatus == model.LicenseStatusActive && (customer.LicenseExpiresDate == nil || customer.LicenseExpiresDate.After(time.Now())) {
		return true, nil
	} else {
		return false, nil

	}

}

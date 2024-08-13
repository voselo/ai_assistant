package handler

import (
	"ai_assistant/config"
	"ai_assistant/internal/customers/model/dto"
	"ai_assistant/internal/customers/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomersHandler struct {
	cfg           *config.Config
	customersRepo repository.CustomersRepository
}

func New(cfg *config.Config, repository repository.CustomersRepository) *CustomersHandler {
	return &CustomersHandler{
		cfg:           cfg,
		customersRepo: repository,
	}
}

// Customer create
//
// @Tags			customers
// @Description		add a new license
// @Router       	/ai/api/v1/customers/create [post]
// @Param        	customer  body      dto.CustomerCreateDTO  true  "customer creation model"
// @Accept       	json
// @Produce      	json
// @Security		ApiKeyAuth
func (handler *CustomersHandler) Create(ctx *gin.Context) {
	var input dto.CustomerCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	customerModel := input.ToModel()
	createdCustomer, err := handler.customersRepo.Create(customerModel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 2"})
		return
	}

	ctx.JSON(http.StatusOK, createdCustomer)
}

// Customer update
//
// @Tags			customers
// @Description		update a license
// @Router       	/ai/api/v1/customers/update/{id} [put]
// @Param			id		path	string				true	"id"
// @Param        	customer  body object  true  "customer update model" example({})
// @Accept       	json
// @Produce      	json
// @Security		ApiKeyAuth
func (handler *CustomersHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input dto.CustomerUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 2"})
		return
	}

	updatedCustomer, err := handler.customersRepo.Update(id, &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 3"})
		return
	}

	ctx.JSON(http.StatusOK, updatedCustomer)

}

// Customer get all
//
// @Tags			customers
// @Description		get all customers
// @Router       	/ai/api/v1/customers/get [get]
// @Accept       	json
// @Security		ApiKeyAuth
func (handler *CustomersHandler) GetAll(ctx *gin.Context) {

	customers, err := handler.customersRepo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 1"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"customers": customers})
}

// Customer get by id
//
// @Tags			customers
// @Description		get customer by id
// @Router       	/ai/api/v1/customers/get/{id} [get]
// @Accept       	json
// @Param        	id path string true "id"
// @Security		ApiKeyAuth
func (handler *CustomersHandler) GetById(ctx *gin.Context) {

	id := ctx.Param("id")
	customer, err := handler.customersRepo.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Err 1"})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

// Customer delete
//
// @Tags			customers
// @Description		delete customer by id
// @Router       	/ai/api/v1/customers/delete/{id} [delete]
// @Accept       	json
// @Param        	id path string true "id"
// @Security		ApiKeyAuth
func (handler *CustomersHandler) Delete(ctx *gin.Context) {

	id := ctx.Param("id")

	err := handler.customersRepo.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 1"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

package handler

import (
	"messages_handler/internal/customers/model/dto"
	"messages_handler/internal/customers/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomersHandler struct {
	customersRepo repository.CustomersRepository
}

func New(repository repository.CustomersRepository) *CustomersHandler {
	return &CustomersHandler{
		customersRepo: repository,
	}
}

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

	ctx.JSON(http.StatusCreated, gin.H{"customer": createdCustomer})
}

func (handler *CustomersHandler) Update(ctx *gin.Context) {
	var input dto.CustomerUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	_, err := uuid.Parse(input.Id.String())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 2"})
		return
	}

	updatedCustomer, err := handler.customersRepo.Update(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 3"})
		return
	}

	ctx.JSON(http.StatusOK, updatedCustomer)

}

func (handler *CustomersHandler) Get(ctx *gin.Context) {
	id := ctx.Query("id")

	if id != "" {
		handler.GetById(ctx, id)
	} else {
		handler.GetAll(ctx)
	}
}

func (handler *CustomersHandler) GetById(ctx *gin.Context, id string) {
	customer, err := handler.customersRepo.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (handler *CustomersHandler) GetAll(ctx *gin.Context) {
	customers, err := handler.customersRepo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve customers"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"customers": customers})
}

func (handler *CustomersHandler) Delete(ctx *gin.Context) {
	id := ctx.Query("id")

	err := handler.customersRepo.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 1"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

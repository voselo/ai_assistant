package handler

import (
	"messages_handler/config"
	"messages_handler/internal/customers/model/dto"
	"messages_handler/internal/customers/repository"
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

func (handler *CustomersHandler) Create(ctx *gin.Context) {
	key := ctx.GetHeader("x-token")

	if key != handler.cfg.ApiKey {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	var input dto.CustomerCreateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 2"})
		return
	}

	customerModel := input.ToModel()
	createdCustomer, err := handler.customersRepo.Create(customerModel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 3"})
		return
	}

	ctx.JSON(http.StatusCreated, createdCustomer)
}

func (handler *CustomersHandler) Update(ctx *gin.Context) {
	key := ctx.GetHeader("x-token")

	if key != handler.cfg.ApiKey {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	var input dto.CustomerUpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 2"})
		return
	}

	_, err := uuid.Parse(input.Id.String())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 3"})
		return
	}

	updatedCustomer, err := handler.customersRepo.Update(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 4"})
		return
	}

	ctx.JSON(http.StatusOK, updatedCustomer)

}

func (handler *CustomersHandler) Get(ctx *gin.Context) {
	key := ctx.GetHeader("x-token")

	if key != handler.cfg.ApiKey {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	id := ctx.Query("id")

	if id != "" {
		handler.GetById(ctx, id)
	} else {
		handler.GetAll(ctx)
	}
}

func (handler *CustomersHandler) GetById(ctx *gin.Context, id string) {
	key := ctx.GetHeader("x-token")

	if key != handler.cfg.ApiKey {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	customer, err := handler.customersRepo.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (handler *CustomersHandler) GetAll(ctx *gin.Context) {
	key := ctx.GetHeader("x-token")

	if key != handler.cfg.ApiKey {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	customers, err := handler.customersRepo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 2"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"customers": customers})
}

func (handler *CustomersHandler) Delete(ctx *gin.Context) {
	key := ctx.GetHeader("x-token")

	if key != handler.cfg.ApiKey {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	id := ctx.Query("id")

	err := handler.customersRepo.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 2"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

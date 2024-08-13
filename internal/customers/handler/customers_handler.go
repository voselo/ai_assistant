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

// Customer create
//
//	@Tags			customers
//	@Description	generates a license
//	@Accept       	json
//	@Produce      	json
//	@Param		 	customer	body	dto.CustomerCreateDTO	true	"create manager"
//	@Accept			json
//	@Produce		json
//	@Success		200  {object}  model.CustomerModel
//	@Router       	/ai/api/v1/customers/create [post]
//	@Security		ApiKeyAuth
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

	ctx.JSON(http.StatusCreated, createdCustomer)
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Err 1"})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (handler *CustomersHandler) GetAll(ctx *gin.Context) {

	customers, err := handler.customersRepo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Err 1"})
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

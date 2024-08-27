package handler

import (
	"authentication/business"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

type getStockByIdController struct {
	service *business.GetStockByIdService
}

func NewGetStockByIdController(service *business.GetStockByIdService) *getStockByIdController {
	return &getStockByIdController{
		service: service,
	}
}

// HandleGetstocks handles the stocks request.
// @Summary Handle get stock by id request
// @Description Handle get stock by id request
// @Tags Get stock by id
// @Produce  json
// @Success 200 {object} string "Returns stocks data"
// @Failure 404 {object} string "Not found"
// @Router /v1/stock/{id} [get]
func (controller *getStockByIdController) HandleGetStockById(context *gin.Context) {
	idStr := context.Param("id")
	idInt, err := strconv.Atoi(idStr)
	id := uint(idInt)
	stock,err := controller.service.GetStockById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message" : "Stock not found"})
		return
	}
	context.IndentedJSON(http.StatusOK,stock)
}

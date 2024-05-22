package handler

import (
	"net/http"
	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
	"stock_broker_application/src/utils/validations"
	"watchlist/business"
	"watchlist/commons/constants"
	"watchlist/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type deleteWatchListController struct {
	service *business.DeleteWatchListService
}

func NewDeleteWatchListDeleteController(service *business.DeleteWatchListService) *deleteWatchListController {
	return &deleteWatchListController{
		service: service,
	}

}

// @Summary Delete a watchlist
// @Description Delete a watchlist with the provided details
// @Tags DeleteWatchlist
// @Accept json
// @Produce json
// @Security JWT
// @Param watchlist body models.WatchlistDeleteModel true "Watchlist details"
// @Success 200 {string} string "Watchlist deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /v1/delete-watchlist [delete]
func (controller *deleteWatchListController) DeleteWatchList(ctx *gin.Context) {
	var watchlist models.DeleteWatchlistRequest
	if err := ctx.ShouldBindJSON(&watchlist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{genericConstants.GenericJSONErrorMessage: err.Error()})
		return
	}
	if err := validations.GetCustomValidator(ctx.Request.Context()).Struct(watchlist); err != nil {
		validationErrors := validations.FormatValidationErrors(ctx.Request.Context(), err.(validator.ValidationErrors))
		ctx.JSON(http.StatusBadRequest, gin.H{genericConstants.GenericJSONErrorMessage: err.Error(),
			genericConstants.GenericValidationError: validationErrors,
		})
	}
	if err := controller.service.DeleteWatchList(&watchlist, ctx); err != nil {
		if err.Error() == constants.NoWatchlistError {
			ctx.JSON(http.StatusNoContent, gin.H{genericConstants.GenericJSONErrorMessage: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{genericConstants.GenericJSONErrorMessage: err.Error()})
		return
	}
	utils.SendStatusOkSuccess(ctx, constants.DeletedWatchlistSuccessMessage)
}

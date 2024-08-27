package business

import (
	"authentication/repositories"
	dbModels "stock_broker_application/src/models"
)

type GetStockByIdService struct {
	stockRepository repositories.GetStockByIdRepository
}

func NewGetStockByIdService(stockRepository repositories.GetStockByIdRepository) *GetStockByIdService {
	return &GetStockByIdService{
		stockRepository: stockRepository,
	}
}

func (service *GetStockByIdService) GetStockById(id uint) (dbModels.Stocks, error){
	var stock dbModels.Stocks
	stock, err := service.stockRepository.GetStockById(id)
	if err != nil {
		return stock,err
	}
	return stock,nil
}

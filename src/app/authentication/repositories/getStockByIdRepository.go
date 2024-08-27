package repositories

import (
	"errors"
	"gorm.io/gorm"
	dbModels "stock_broker_application/src/models"
)

type GetStockByIdRepository interface {
	GetStockById(id uint)(dbModels.Stocks, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *stockRepository {
	return &stockRepository{db: db}
}

func (repository *stockRepository) GetStockById(id uint)(dbModels.Stocks, error) {
	var stock dbModels.Stocks
	if err := repository.db.Where("id = ?", id).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stock, errors.New("Not found")
		}
		return stock, errors.New("Not found")
	}
	return stock, nil
}



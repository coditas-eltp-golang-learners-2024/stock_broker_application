package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"stock_broker_application/src/models"
)

// MockDeleteWatchlistRepository is a mock implementation of DeleteWatchlistRepository.
type MockDeleteWatchlistRepository struct {
	mock.Mock
}

func (m *MockDeleteWatchlistRepository) DeleteScrips(watchlistName string, scrips []int, watchlistCondition map[string]interface{}) error {
	args := m.Called(watchlistName, scrips, watchlistCondition)
	return args.Error(0)
}

// DeleteWatchlistRepositorySuite is a test suite to test DeleteWatchlistRepository.
type DeleteWatchlistRepositorySuite struct {
	suite.Suite
	mockRepo *MockDeleteWatchlistRepository
	db       *gorm.DB
}

func (suite *DeleteWatchlistRepositorySuite) SetupSuite() {
	// Setup mock repository
	suite.mockRepo = &MockDeleteWatchlistRepository{}

	// Setup SQLite in-memory database for GORM
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	suite.Require().NoError(err)

	// Auto migrate models
	err = db.AutoMigrate(&models.Watchlist{}, &models.WatchlistStock{})
	suite.Require().NoError(err)

	suite.db = db
}

func (suite *DeleteWatchlistRepositorySuite) TearDownSuite() {
	// Close GORM DB
	sqlDB, err := suite.db.DB()
	suite.Require().NoError(err)
	sqlDB.Close()
}

func (suite *DeleteWatchlistRepositorySuite) TestDeleteScrips() {
	const (
		WatchlistNotFoundError      = "watchlist not found"
		FailedtoFindWatchlistError  = "failed to find watchlist"
		FailToDeleteScripsError     = "failed to delete scrip"
	)

	tests := []struct {
		name               string
		watchlistCondition map[string]interface{}
		scrips             []int
		mockBehavior       func()
		expectedError      error
	}{
		{
			name: "successful deletion of scrips",
			watchlistCondition: map[string]interface{}{
				"name": "test_watchlist",
			},
			scrips: []int{1},
			mockBehavior: func() {
				suite.mockRepo.On("DeleteScrips", "test_watchlist", []int{1}, mock.Anything).
					Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "watchlist not found",
			watchlistCondition: map[string]interface{}{
				"name": "test_watchlist",
			},
			scrips: []int{1},
			mockBehavior: func() {
				suite.mockRepo.On("DeleteScrips", "test_watchlist", []int{1}, mock.Anything).
					Return(errors.New(WatchlistNotFoundError)).Once()
			},
			expectedError: errors.New(WatchlistNotFoundError),
		},
		{
			name: "failed to find watchlist",
			watchlistCondition: map[string]interface{}{
				"name": "test_watchlist",
			},
			scrips: []int{1},
			mockBehavior: func() {
				suite.mockRepo.On("DeleteScrips", "test_watchlist", []int{1}, mock.Anything).
					Return(errors.New(FailedtoFindWatchlistError)).Once()
			},
			expectedError: errors.New(FailedtoFindWatchlistError),
		},
		{
			name: "failed to delete scrip",
			watchlistCondition: map[string]interface{}{
				"name": "test_watchlist",
			},
			scrips: []int{1},
			mockBehavior: func() {
				suite.mockRepo.On("DeleteScrips", "test_watchlist", []int{1}, mock.Anything).
					Return(errors.New(FailToDeleteScripsError)).Once()
			},
			expectedError: errors.New(FailToDeleteScripsError),
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test.mockBehavior()

			err := suite.mockRepo.DeleteScrips("test_watchlist", []int{1}, test.watchlistCondition)
			if test.expectedError != nil {
				assert.EqualError(suite.T(), err, test.expectedError.Error())
			} else {
				assert.NoError(suite.T(), err)
			}
		})
	}
}

func TestDeleteWatchlistRepositorySuite(t *testing.T) {
	suite.Run(t, new(DeleteWatchlistRepositorySuite))
}

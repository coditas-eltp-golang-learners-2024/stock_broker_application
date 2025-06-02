package tests

import (
	"authentication/commons/constants"
	"authentication/repositories"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func setupTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gdb, mock, nil
}

func checkExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	assert.NoError(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
}

func TestCheckOtp(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)
	repo := repositories.NewUserRepository(db)

	tests := []struct {
		name      string
		userID    uint16
		otp       uint16
		setupMock func()
		expected  bool
		expectErr bool
	}{
		{
			name:   "valid OTP",
			userID: 1,
			otp:    1234,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1 AND otp = $2`)).
					WithArgs(1, 1234).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:  true,
			expectErr: false,
		},
		{
			name:   "invalid OTP",
			userID: 1,
			otp:    1234,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1 AND otp = $2`)).
					WithArgs(1, 1234).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:  false,
			expectErr: false,
		},
		{
			name:   "query error",
			userID: 1,
			otp:    1234,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1 AND otp = $2`)).
					WithArgs(1, 1234).
					WillReturnError(errors.New("query error"))
			},
			expected:  false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			valid, err := repo.CheckOtp(tt.userID, tt.otp)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, constants.ErrorUpdateUserToken+"query error", err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, valid)
			checkExpectations(t, mock)
		})
	}
}

func TestCheckUserExists(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)
	repo := repositories.NewUserRepository(db)

	tests := []struct {
		name      string
		userID    uint16
		setupMock func()
		expected  bool
		expectErr bool
	}{
		{
			name:   "user exists",
			userID: 1,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1`)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:  true,
			expectErr: false,
		},
		{
			name:   "user does not exist",
			userID: 1,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1`)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:  false,
			expectErr: false,
		},
		{
			name:   "query error",
			userID: 1,
			setupMock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE id = $1`)).
					WithArgs(1).
					WillReturnError(errors.New("query error"))
			},
			expected:  false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			exists, err := repo.CheckUserExists(tt.userID)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, exists)
			checkExpectations(t, mock)
		})
	}
}

func TestUpdateUserToken(t *testing.T) {
	db, mock, err := setupTestDB()
	assert.NoError(t, err)
	repo := repositories.NewUserRepository(db)

	tests := []struct {
		name      string
		userID    uint16
		token     string
		setupMock func()
		expectErr bool
	}{
		{
			name:   "update_token_success",
			userID: 1,
			token:  "new_token",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "token"=$1 WHERE id = $2`)).
					WithArgs("new_token", uint16(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:   "update_token_error",
			userID: 1,
			token:  "new_token",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "token"=$1 WHERE id = $2`)).
					WithArgs("new_token", uint16(1)).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.UpdateUserToken(tt.userID, tt.token)
			if tt.expectErr {
				assert.EqualError(t, err, constants.ErrorUpdateUserToken+"update error")
			} else {
				assert.NoError(t, err)
			}
			checkExpectations(t, mock)
		})
	}
}

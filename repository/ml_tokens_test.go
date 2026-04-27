package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
)

func setupMLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *mlTokensDAO) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}

	dao := &mlTokensDAO{db: db}

	return db, mock, dao
}

func TestGetByUserID_Success(t *testing.T) {
	db, mock, dao := setupMLMock(t)
	defer db.Close()

	now := time.Now()
	nickname := "gabriel"

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "nickname",
		"access_token", "refresh_token",
		"expires_at", "created_at", "updated_at",
	}).AddRow(
		1, 123, nickname,
		"access", "refresh",
		now, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(int64(123)).
		WillReturnRows(rows)

	token, err := dao.GetByUserID(123)

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}

	if token == nil {
		t.Errorf("token não deveria ser nil")
	}

	if token.UserID != 123 {
		t.Errorf("user_id incorreto")
	}
}

func TestGetByUserID_NotFound(t *testing.T) {
	db, mock, dao := setupMLMock(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(int64(123)).
		WillReturnError(sql.ErrNoRows)

	token, err := dao.GetByUserID(123)

	if err != nil {
		t.Errorf("não deveria retornar erro")
	}

	if token != nil {
		t.Errorf("token deveria ser nil")
	}
}

func TestUpsert_Success(t *testing.T) {
	db, mock, dao := setupMLMock(t)
	defer db.Close()

	token := &model.MLToken{
		UserID:       123,
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO ml_tokens")).
		WithArgs(
			token.UserID,
			nil,
			token.AccessToken,
			token.RefreshToken,
			token.ExpiresAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Upsert(token)

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}

func TestUpdateTokens_Success(t *testing.T) {
	db, mock, dao := setupMLMock(t)
	defer db.Close()

	expires := time.Now()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE ml_tokens")).
		WithArgs(
			"access",
			"refresh",
			expires,
			int64(123),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.UpdateTokens(123, "access", "refresh", expires)

	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}

func TestGetValidToken_Valid(t *testing.T) {
	db, mock, dao := setupMLMock(t)
	defer db.Close()

	now := time.Now().Add(1 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "nickname",
		"access_token", "refresh_token",
		"expires_at", "created_at", "updated_at",
	}).AddRow(
		1, 123, nil,
		"access", "refresh",
		now, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(int64(123)).
		WillReturnRows(rows)

	token, err := dao.GetValidToken(123)

	if err != nil {
		t.Errorf("não deveria dar erro")
	}

	if token == nil {
		t.Errorf("token não deveria ser nil")
	}
}

func TestGetValidToken_Expired(t *testing.T) {
	db, mock, dao := setupMLMock(t)
	defer db.Close()

	now := time.Now().Add(-1 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "nickname",
		"access_token", "refresh_token",
		"expires_at", "created_at", "updated_at",
	}).AddRow(
		1, 123, nil,
		"access", "refresh",
		now, now, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(int64(123)).
		WillReturnRows(rows)

	token, err := dao.GetValidToken(123)

	if err == nil {
		t.Errorf("deveria retornar erro de expirado")
	}

	if token == nil {
		t.Errorf("token ainda deveria ser retornado")
	}
}

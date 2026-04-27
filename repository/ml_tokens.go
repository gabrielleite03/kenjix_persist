package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type MLTokensDAO interface {
	GetByUserID(userID int64) (*model.MLToken, error)
	Upsert(token *model.MLToken) error
	UpdateTokens(userID int64, accessToken string, refreshToken string, expiresAt time.Time) error
	GetValidToken(userID int64) (*model.MLToken, error)
}

type mlTokensDAO struct {
	db *sql.DB
}

func NewMLTokensDAO() MLTokensDAO {
	return &mlTokensDAO{db: config.NewDatabaseConfig().DB}
}

func (r *mlTokensDAO) GetByUserID(userID int64) (*model.MLToken, error) {
	query := `
		SELECT 
			id, user_id, nickname,
			access_token, refresh_token,
			expires_at, created_at, updated_at
		FROM ml_tokens
		WHERE user_id = ?
		LIMIT 1
	`

	row := r.db.QueryRow(query, userID)

	var token model.MLToken
	var nickname sql.NullString

	err := row.Scan(
		&token.ID,
		&token.UserID,
		&nickname,
		&token.AccessToken,
		&token.RefreshToken,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if nickname.Valid {
		token.Nickname = &nickname.String
	}

	return &token, nil
}

func (r *mlTokensDAO) Upsert(token *model.MLToken) error {
	query := `
		INSERT INTO ml_tokens (
			user_id,
			nickname,
			access_token,
			refresh_token,
			expires_at
		) VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			nickname = VALUES(nickname),
			access_token = VALUES(access_token),
			refresh_token = VALUES(refresh_token),
			expires_at = VALUES(expires_at),
			updated_at = CURRENT_TIMESTAMP
	`

	var nickname interface{}
	if token.Nickname != nil {
		nickname = *token.Nickname
	}

	_, err := r.db.Exec(
		query,
		token.UserID,
		nickname,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresAt,
	)

	return err
}

func (r *mlTokensDAO) UpdateTokens(
	userID int64,
	accessToken string,
	refreshToken string,
	expiresAt time.Time,
) error {

	query := `
		UPDATE ml_tokens
		SET 
			access_token = ?,
			refresh_token = ?,
			expires_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = ?
	`

	_, err := r.db.Exec(
		query,
		accessToken,
		refreshToken,
		expiresAt,
		userID,
	)

	return err
}

func (r *mlTokensDAO) GetValidToken(userID int64) (*model.MLToken, error) {
	token, err := r.GetByUserID(userID)
	if err != nil || token == nil {
		return token, err
	}

	if time.Now().After(token.ExpiresAt) {
		return token, errors.New("token expirado")
	}

	return token, nil
}

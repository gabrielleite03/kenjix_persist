package repository

import (
	"context"
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type NFeIssuerDAO interface {
	Create(ctx context.Context, issuer *model.NFeIssuer) error
	Update(ctx context.Context, issuer *model.NFeIssuer) error
	FindByID(ctx context.Context, id int64) (*model.NFeIssuer, error)
	FindAll(ctx context.Context) ([]model.NFeIssuer, error)
	Delete(ctx context.Context, id int64) error
}

type nFeIssuerDAO struct {
	db *sql.DB
}

func NewNFeIssuerDAO() NFeIssuerDAO {
	return &nFeIssuerDAO{db: config.NewDatabaseConfig().DB}
}

func (r *nFeIssuerDAO) Create(ctx context.Context, issuer *model.NFeIssuer) error {
	query := `
		INSERT INTO nfe_issuers (
			cnpj,
			cpf,
			razao_social,
			nome_fantasia,
			inscricao_estadual,
			crt,
			logradouro,
			numero,
			complemento,
			bairro,
			codigo_mun,
			municipio,
			uf,
			cep,
			codigo_pais,
			pais,
			telefone,
			active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		stringValue(issuer.CNPJ),
		stringValue(issuer.CPF),
		issuer.RazaoSocial,
		stringValue(issuer.NomeFantasia),
		stringValue(issuer.InscricaoEstadual),
		issuer.CRT,
		issuer.Logradouro,
		issuer.Numero,
		stringValue(issuer.Complemento),
		issuer.Bairro,
		issuer.CodigoMun,
		issuer.Municipio,
		issuer.UF,
		issuer.CEP,
		issuer.CodigoPais,
		issuer.Pais,
		stringValue(issuer.Telefone),
		issuer.Active,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	issuer.ID = id
	return nil
}

func (r *nFeIssuerDAO) Update(ctx context.Context, issuer *model.NFeIssuer) error {
	query := `
		UPDATE nfe_issuers SET
			cnpj = ?,
			cpf = ?,
			razao_social = ?,
			nome_fantasia = ?,
			inscricao_estadual = ?,
			crt = ?,
			logradouro = ?,
			numero = ?,
			complemento = ?,
			bairro = ?,
			codigo_mun = ?,
			municipio = ?,
			uf = ?,
			cep = ?,
			codigo_pais = ?,
			pais = ?,
			telefone = ?,
			active = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		stringValue(issuer.CNPJ),
		stringValue(issuer.CPF),
		issuer.RazaoSocial,
		stringValue(issuer.NomeFantasia),
		stringValue(issuer.InscricaoEstadual),
		issuer.CRT,
		issuer.Logradouro,
		issuer.Numero,
		stringValue(issuer.Complemento),
		issuer.Bairro,
		issuer.CodigoMun,
		issuer.Municipio,
		issuer.UF,
		issuer.CEP,
		issuer.CodigoPais,
		issuer.Pais,
		stringValue(issuer.Telefone),
		issuer.Active,
		issuer.ID,
	)

	return err
}

func (r *nFeIssuerDAO) FindByID(ctx context.Context, id int64) (*model.NFeIssuer, error) {
	query := `
		SELECT
			id,
			cnpj,
			cpf,
			razao_social,
			nome_fantasia,
			inscricao_estadual,
			crt,
			logradouro,
			numero,
			complemento,
			bairro,
			codigo_mun,
			municipio,
			uf,
			cep,
			codigo_pais,
			pais,
			telefone,
			active,
			created_at,
			updated_at
		FROM nfe_issuers
		WHERE id = ? AND active = 1
	`

	var issuer model.NFeIssuer
	scanValues := newNFeIssuerScanValues(&issuer)
	err := r.db.QueryRowContext(ctx, query, id).Scan(scanValues.fields()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	scanValues.apply()

	return &issuer, nil
}

func (r *nFeIssuerDAO) FindAll(ctx context.Context) ([]model.NFeIssuer, error) {
	query := `
		SELECT
			id,
			cnpj,
			cpf,
			razao_social,
			nome_fantasia,
			inscricao_estadual,
			crt,
			logradouro,
			numero,
			complemento,
			bairro,
			codigo_mun,
			municipio,
			uf,
			cep,
			codigo_pais,
			pais,
			telefone,
			active,
			created_at,
			updated_at
		FROM nfe_issuers
		WHERE active = 1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.NFeIssuer
	for rows.Next() {
		var issuer model.NFeIssuer
		scanValues := newNFeIssuerScanValues(&issuer)
		if err := rows.Scan(scanValues.fields()...); err != nil {
			return nil, err
		}
		scanValues.apply()

		list = append(list, issuer)
	}

	return list, nil
}

func (r *nFeIssuerDAO) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE nfe_issuers
		SET active = 0
		WHERE id = ? AND active = 1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

type nFeIssuerScanValues struct {
	issuer            *model.NFeIssuer
	cnpj              sql.NullString
	cpf               sql.NullString
	nomeFantasia      sql.NullString
	inscricaoEstadual sql.NullString
	complemento       sql.NullString
	telefone          sql.NullString
}

func newNFeIssuerScanValues(issuer *model.NFeIssuer) *nFeIssuerScanValues {
	return &nFeIssuerScanValues{issuer: issuer}
}

func (s *nFeIssuerScanValues) fields() []interface{} {
	return []interface{}{
		&s.issuer.ID,
		&s.cnpj,
		&s.cpf,
		&s.issuer.RazaoSocial,
		&s.nomeFantasia,
		&s.inscricaoEstadual,
		&s.issuer.CRT,
		&s.issuer.Logradouro,
		&s.issuer.Numero,
		&s.complemento,
		&s.issuer.Bairro,
		&s.issuer.CodigoMun,
		&s.issuer.Municipio,
		&s.issuer.UF,
		&s.issuer.CEP,
		&s.issuer.CodigoPais,
		&s.issuer.Pais,
		&s.telefone,
		&s.issuer.Active,
		&s.issuer.CreatedAt,
		&s.issuer.UpdatedAt,
	}
}

func (s *nFeIssuerScanValues) apply() {
	if s.cnpj.Valid {
		s.issuer.CNPJ = &s.cnpj.String
	}
	if s.cpf.Valid {
		s.issuer.CPF = &s.cpf.String
	}
	if s.nomeFantasia.Valid {
		s.issuer.NomeFantasia = &s.nomeFantasia.String
	}
	if s.inscricaoEstadual.Valid {
		s.issuer.InscricaoEstadual = &s.inscricaoEstadual.String
	}
	if s.complemento.Valid {
		s.issuer.Complemento = &s.complemento.String
	}
	if s.telefone.Valid {
		s.issuer.Telefone = &s.telefone.String
	}
}

func stringValue(value *string) interface{} {
	if value == nil {
		return nil
	}

	return *value
}

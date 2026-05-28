package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
)

func setupNFeIssuerMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, NFeIssuerDAO) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dao := &nFeIssuerDAO{db: db}

	return db, mock, dao
}

func TestCreateNFeIssuer(t *testing.T) {
	db, mock, dao := setupNFeIssuerMock(t)
	defer db.Close()

	issuer := newNFeIssuer()

	mock.ExpectExec(regexp.QuoteMeta(`
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
	`)).
		WithArgs(
			*issuer.CNPJ,
			nil,
			issuer.RazaoSocial,
			*issuer.NomeFantasia,
			*issuer.InscricaoEstadual,
			issuer.CRT,
			issuer.Logradouro,
			issuer.Numero,
			nil,
			issuer.Bairro,
			issuer.CodigoMun,
			issuer.Municipio,
			issuer.UF,
			issuer.CEP,
			issuer.CodigoPais,
			issuer.Pais,
			*issuer.Telefone,
			issuer.Active,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Create(context.Background(), issuer)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), issuer.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateNFeIssuer(t *testing.T) {
	db, mock, dao := setupNFeIssuerMock(t)
	defer db.Close()

	issuer := newNFeIssuer()
	issuer.ID = 1

	mock.ExpectExec(regexp.QuoteMeta(`
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
	`)).
		WithArgs(
			*issuer.CNPJ,
			nil,
			issuer.RazaoSocial,
			*issuer.NomeFantasia,
			*issuer.InscricaoEstadual,
			issuer.CRT,
			issuer.Logradouro,
			issuer.Numero,
			nil,
			issuer.Bairro,
			issuer.CodigoMun,
			issuer.Municipio,
			issuer.UF,
			issuer.CEP,
			issuer.CodigoPais,
			issuer.Pais,
			*issuer.Telefone,
			issuer.Active,
			issuer.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Update(context.Background(), issuer)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindNFeIssuerByID(t *testing.T) {
	db, mock, dao := setupNFeIssuerMock(t)
	defer db.Close()

	now := time.Now()
	rows := nFeIssuerRows().
		AddRow(
			1,
			"12345678000195",
			nil,
			"Empresa Teste LTDA",
			"Empresa Teste",
			"123456789",
			"3",
			"Rua Central",
			"100",
			nil,
			"Centro",
			"3550308",
			"Sao Paulo",
			"SP",
			"01001000",
			"1058",
			"Brasil",
			"11999999999",
			true,
			now,
			now,
		)

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	result, err := dao.FindByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "12345678000195", *result.CNPJ)
	assert.Nil(t, result.CPF)
	assert.Equal(t, "Empresa Teste LTDA", result.RazaoSocial)
	assert.Nil(t, result.Complemento)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllNFeIssuers(t *testing.T) {
	db, mock, dao := setupNFeIssuerMock(t)
	defer db.Close()

	now := time.Now()
	rows := nFeIssuerRows().
		AddRow(1, "12345678000195", nil, "Empresa A", nil, nil, "3", "Rua A", "10", nil, "Centro", "3550308", "Sao Paulo", "SP", "01001000", "1058", "Brasil", nil, true, now, now).
		AddRow(2, nil, "12345678901", "Empresa B", nil, nil, "1", "Rua B", "20", nil, "Centro", "3304557", "Rio de Janeiro", "RJ", "20040002", "1058", "Brasil", nil, true, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`
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
	`)).
		WillReturnRows(rows)

	result, err := dao.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "12345678901", *result[1].CPF)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteNFeIssuer(t *testing.T) {
	db, mock, dao := setupNFeIssuerMock(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE nfe_issuers
		SET active = 0
		WHERE id = ? AND active = 1
	`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func newNFeIssuer() *model.NFeIssuer {
	return &model.NFeIssuer{
		CNPJ:              nFeIssuerString("12345678000195"),
		RazaoSocial:       "Empresa Teste LTDA",
		NomeFantasia:      nFeIssuerString("Empresa Teste"),
		InscricaoEstadual: nFeIssuerString("123456789"),
		CRT:               "3",
		Logradouro:        "Rua Central",
		Numero:            "100",
		Bairro:            "Centro",
		CodigoMun:         "3550308",
		Municipio:         "Sao Paulo",
		UF:                "SP",
		CEP:               "01001000",
		CodigoPais:        "1058",
		Pais:              "Brasil",
		Telefone:          nFeIssuerString("11999999999"),
		Active:            true,
	}
}

func nFeIssuerString(value string) *string {
	return &value
}

func nFeIssuerRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id",
		"cnpj",
		"cpf",
		"razao_social",
		"nome_fantasia",
		"inscricao_estadual",
		"crt",
		"logradouro",
		"numero",
		"complemento",
		"bairro",
		"codigo_mun",
		"municipio",
		"uf",
		"cep",
		"codigo_pais",
		"pais",
		"telefone",
		"active",
		"created_at",
		"updated_at",
	})
}

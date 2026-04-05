package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type ExpenseDAO interface {
	FindAll() ([]model.Expense, error)
	FindByID(id int64) (*model.Expense, error)
	Create(expense *model.Expense) error
	Update(expense *model.Expense) error
	Delete(id int64) error
	AddAttachment(expenseID int64, url string) error
	RemoveAttachment(id int64) error
	FindAllExpenseCategories() ([]model.ExpenseCategory, error)
}

type expenseDAO struct {
	db *sql.DB
}

func NewExpenseDAO() ExpenseDAO {
	return &expenseDAO{db: config.NewDatabaseConfig().DB}
}

func (d *expenseDAO) FindAll() ([]model.Expense, error) {
	query := `
	SELECT 
		e.id, e.description, e.category_id, e.amount, e.date, e.status,
		e.created_at, e.updated_at,
		c.id, c.name, c.created_at, c.updated_at
	FROM expenses e
	LEFT JOIN expenses_category c ON c.id = e.category_id
	ORDER BY e.date DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []model.Expense

	for rows.Next() {
		var exp model.Expense
		var cat model.ExpenseCategory

		err := rows.Scan(
			&exp.ID,
			&exp.Description,
			&exp.CategoryID,
			&exp.Amount,
			&exp.Date,
			&exp.Status,
			&exp.CreatedAt,
			&exp.UpdatedAt,
			&cat.ID,
			&cat.Name,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		exp.Category = &cat

		attachments, err := d.getAttachments(exp.ID)
		if err != nil {
			return nil, err
		}
		exp.Attachments = attachments

		expenses = append(expenses, exp)
	}

	return expenses, nil
}

func (d *expenseDAO) FindByID(id int64) (*model.Expense, error) {
	query := `
	SELECT 
		e.id, e.description, e.category_id, e.amount, e.date, e.status,
		e.created_at, e.updated_at,
		c.id, c.name, c.created_at, c.updated_at
	FROM expenses e
	LEFT JOIN expenses_category c ON c.id = e.category_id
	WHERE e.id = ?
	`

	row := d.db.QueryRow(query, id)

	var exp model.Expense
	var cat model.ExpenseCategory

	err := row.Scan(
		&exp.ID,
		&exp.Description,
		&exp.CategoryID,
		&exp.Amount,
		&exp.Date,
		&exp.Status,
		&exp.CreatedAt,
		&exp.UpdatedAt,
		&cat.ID,
		&cat.Name,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	exp.Category = &cat

	attachments, err := d.getAttachments(id)
	if err != nil {
		return nil, err
	}
	exp.Attachments = attachments

	return &exp, nil
}

func (d *expenseDAO) Create(exp *model.Expense) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO expenses
	(description, category_id, amount, date, status)
	VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		query,
		exp.Description,
		exp.CategoryID,
		exp.Amount,
		exp.Date,
		exp.Status,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	exp.ID = id

	// salvar attachments
	for _, att := range exp.Attachments {
		_, err := tx.Exec(
			"INSERT INTO expense_attachments (expense_id, url) VALUES (?, ?)",
			exp.ID,
			att.URL,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *expenseDAO) Update(exp *model.Expense) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := `
	UPDATE expenses SET
	description = ?,
	category_id = ?,
	amount = ?,
	date = ?,
	status = ?
	WHERE id = ?
	`

	_, err = tx.Exec(
		query,
		exp.Description,
		exp.CategoryID,
		exp.Amount,
		exp.Date,
		exp.Status,
		exp.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// remove attachments antigos
	_, err = tx.Exec(
		"DELETE FROM expense_attachments WHERE expense_id = ?",
		exp.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// reinsere attachments
	for _, att := range exp.Attachments {
		_, err := tx.Exec(
			"INSERT INTO expense_attachments (expense_id, url) VALUES (?, ?)",
			exp.ID,
			att.URL,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *expenseDAO) Delete(id int64) error {
	_, err := d.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}

func (d *expenseDAO) getAttachments(expenseID int64) ([]model.ExpenseAttachment, error) {
	rows, err := d.db.Query(
		"SELECT id, expense_id, url, created_at FROM expense_attachments WHERE expense_id = ?",
		expenseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []model.ExpenseAttachment

	for rows.Next() {
		var att model.ExpenseAttachment
		err := rows.Scan(
			&att.ID,
			&att.ExpenseID,
			&att.URL,
			&att.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, att)
	}

	return attachments, nil
}

func (d *expenseDAO) AddAttachment(expenseID int64, url string) error {
	_, err := d.db.Exec(
		"INSERT INTO expense_attachments (expense_id, url) VALUES (?, ?)",
		expenseID,
		url,
	)
	return err
}

func (d *expenseDAO) RemoveAttachment(id int64) error {
	_, err := d.db.Exec(
		"DELETE FROM expense_attachments WHERE id = ?",
		id,
	)
	return err
}

func (d *expenseDAO) FindAllExpenseCategories() ([]model.ExpenseCategory, error) {
	rows, err := d.db.Query(`
		SELECT id, name, created_at, updated_at
		FROM expenses_category
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.ExpenseCategory

	for rows.Next() {
		var c model.ExpenseCategory
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

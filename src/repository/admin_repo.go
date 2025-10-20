package repository

import (
	"database/sql"
	"xyz_backend/src/models"
)

type AdminRepository interface {
	Create(admin *models.Admin) error
	FindByEmail(email string) (*models.Admin, error)
	FindByID(id int64) (*models.Admin, error)
}

type adminRepo struct{ db *sql.DB }

func NewAdminRepo(db *sql.DB) AdminRepository { return &adminRepo{db} }

func (r *adminRepo) Create(a *models.Admin) error {
	stmt := `INSERT INTO admins (name, email, password_hash, role) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(stmt, a.Name, a.Email, a.PasswordHash, a.Role)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	a.ID = id
	return nil
}

func (r *adminRepo) FindByEmail(email string) (*models.Admin, error) {
	a := &models.Admin{}
	row := r.db.QueryRow(`SELECT id,name,email,password_hash,role,created_at,updated_at,deleted_at FROM admins WHERE email=? AND deleted_at IS NULL`, email)
	var deleted sql.NullTime
	if err := row.Scan(&a.ID, &a.Name, &a.Email, &a.PasswordHash, &a.Role, &a.CreatedAt, &a.UpdatedAt, &deleted); err != nil {
		return nil, err
	}
	if deleted.Valid {
		t := deleted.Time
		a.DeletedAt = &t
	}
	return a, nil
}

func (r *adminRepo) FindByID(id int64) (*models.Admin, error) {
	a := &models.Admin{}
	var deleted sql.NullTime
	row := r.db.QueryRow(`SELECT id,name,email,password_hash,role,created_at,updated_at,deleted_at FROM admins WHERE id=? AND deleted_at IS NULL`, id)
	if err := row.Scan(&a.ID, &a.Name, &a.Email, &a.PasswordHash, &a.Role, &a.CreatedAt, &a.UpdatedAt, &deleted); err != nil {
		return nil, err
	}
	if deleted.Valid {
		t := deleted.Time
		a.DeletedAt = &t
	}
	return a, nil
}

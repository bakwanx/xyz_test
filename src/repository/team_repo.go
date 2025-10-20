package repository

import (
	"database/sql"
	"time"

	"xyz_backend/src/models"
)

type TeamRepository interface {
	Create(t *models.Team) error
	Update(t *models.Team) error
	SoftDelete(id int64) error
	GetByID(id int64) (*models.Team, error)
	List(offset, limit int) ([]*models.Team, error)
	Count() (int, error)
}

type teamRepo struct{ db *sql.DB }

func NewTeamRepo(db *sql.DB) TeamRepository { return &teamRepo{db} }

func (r *teamRepo) Create(t *models.Team) error {
	res, err := r.db.Exec(`INSERT INTO teams (name,logo_path,founded_year,address,city) VALUES (?,?,?,?,?)`, t.Name, t.LogoPath, t.FoundedYear, t.Address, t.City)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	t.ID = id
	t.CreatedAt = time.Now()
	return nil
}

func (r *teamRepo) Update(t *models.Team) error {
	_, err := r.db.Exec(`UPDATE teams SET name=?,logo_path=?,founded_year=?,address=?,city=?,updated_at=NOW() WHERE id=? AND deleted_at IS NULL`, t.Name, t.LogoPath, t.FoundedYear, t.Address, t.City, t.ID)
	return err
}

func (r *teamRepo) SoftDelete(id int64) error {
	_, err := r.db.Exec(`UPDATE teams SET deleted_at=NOW() WHERE id=?`, id)
	return err
}

func (r *teamRepo) GetByID(id int64) (*models.Team, error) {
	t := &models.Team{}
	var logo sql.NullString
	var deleted sql.NullTime
	row := r.db.QueryRow(`SELECT id,name,logo_path,founded_year,address,city,created_at,updated_at,deleted_at FROM teams WHERE id=?`, id)
	if err := row.Scan(&t.ID, &t.Name, &logo, &t.FoundedYear, &t.Address, &t.City, &t.CreatedAt, &t.UpdatedAt, &deleted); err != nil {
		return nil, err
	}
	if logo.Valid {
		t.LogoPath = &logo.String
	}
	if deleted.Valid {
		dt := deleted.Time
		t.DeletedAt = &dt
	}
	return t, nil
}

func (r *teamRepo) List(offset, limit int) ([]*models.Team, error) {
	rows, err := r.db.Query(`SELECT id,name,logo_path,founded_year,address,city,created_at,updated_at FROM teams WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*models.Team
	for rows.Next() {
		t := &models.Team{}
		var logo sql.NullString
		if err := rows.Scan(&t.ID, &t.Name, &logo, &t.FoundedYear, &t.Address, &t.City, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if logo.Valid {
			t.LogoPath = &logo.String
		}
		res = append(res, t)
	}
	return res, nil
}

func (r *teamRepo) Count() (int, error) {
	var cnt int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM teams WHERE deleted_at IS NULL`).Scan(&cnt)
	return cnt, err
}

package repository

import (
	"database/sql"
	"time"

	"xyz_backend/src/models"
)

type PlayerRepository interface {
	Create(p *models.Player) error
	Update(p *models.Player) error
	SoftDelete(id int64) error
	GetByID(id int64) (*models.Player, error)
	IsJerseyNumberAvailable(teamId int64, jerseyNumber int) (bool, error)
	ListByTeam(teamID int64) ([]*models.Player, error)
}

type playerRepo struct{ db *sql.DB }

func NewPlayerRepo(db *sql.DB) PlayerRepository { return &playerRepo{db} }

func (r *playerRepo) Create(p *models.Player) error {
	res, err := r.db.Exec(`INSERT INTO players (team_id,name,height_cm,weight_kg,position,jersey_number) VALUES (?,?,?,?,?,?)`, p.TeamID, p.Name, p.HeightCm, p.WeightKg, p.Position, p.JerseyNumber)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	p.ID = id
	p.CreatedAt = time.Now()
	return nil
}

func (r *playerRepo) Update(p *models.Player) error {
	_, err := r.db.Exec(`UPDATE players SET team_id=?,name=?,height_cm=?,weight_kg=?,position=?,jersey_number=?,updated_at=NOW() WHERE id=? AND deleted_at IS NULL`, p.TeamID, p.Name, p.HeightCm, p.WeightKg, p.Position, p.JerseyNumber, p.ID)
	return err
}

func (r *playerRepo) SoftDelete(id int64) error {
	_, err := r.db.Exec(`UPDATE players SET deleted_at=NOW() WHERE id=?`, id)
	return err
}

func (r *playerRepo) GetByID(id int64) (*models.Player, error) {
	p := &models.Player{}
	var deleted sql.NullTime
	row := r.db.QueryRow(`SELECT id,team_id,name,height_cm,weight_kg,position,jersey_number,created_at,updated_at,deleted_at FROM players WHERE id=?`, id)
	if err := row.Scan(&p.ID, &p.TeamID, &p.Name, &p.HeightCm, &p.WeightKg, &p.Position, &p.JerseyNumber, &p.CreatedAt, &p.UpdatedAt, &deleted); err != nil {
		return nil, err
	}
	if deleted.Valid {
		t := deleted.Time
		p.DeletedAt = &t
	}
	return p, nil
}

func (r *playerRepo) IsJerseyNumberAvailable(teamId int64, jerseyNumber int) (bool, error) {
	var cnt int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM players WHERE team_id=? AND jersey_number=?`, teamId, jerseyNumber).Scan(&cnt)
	return cnt == 0, err
}

func (r *playerRepo) ListByTeam(teamID int64) ([]*models.Player, error) {
	rows, err := r.db.Query(`SELECT id,team_id,name,height_cm,weight_kg,position,jersey_number,created_at,updated_at FROM players WHERE team_id=? AND deleted_at IS NULL ORDER BY jersey_number ASC`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*models.Player
	for rows.Next() {
		p := &models.Player{}
		if err := rows.Scan(&p.ID, &p.TeamID, &p.Name, &p.HeightCm, &p.WeightKg, &p.Position, &p.JerseyNumber, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

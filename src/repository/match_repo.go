package repository

import (
	"database/sql"
	"time"
	"xyz_backend/src/models"
)

type MatchRepository interface {
	Create(m *models.Match) error
	Update(m *models.Match) error
	GetByID(id int64) (*models.Match, error)
	ListByTeam(teamID int64) ([]*models.Match, error)
	CountWins(teamID int64, asHome bool, upTo time.Time) (int, error)
}

type matchRepo struct{ db *sql.DB }

func NewMatchRepo(db *sql.DB) MatchRepository { return &matchRepo{db} }

func (r *matchRepo) Create(m *models.Match) error {
	res, err := r.db.Exec(`INSERT INTO matches (scheduled_at,home_team_id,away_team_id,status) VALUES (?,?,?,?)`, m.ScheduledAt, m.HomeTeamID, m.AwayTeamID, m.Status)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	m.ID = id
	m.CreatedAt = time.Now()
	return nil
}

func (r *matchRepo) Update(m *models.Match) error {
	_, err := r.db.Exec(`UPDATE matches SET scheduled_at=?,home_team_id=?,away_team_id=?,home_score=?,away_score=?,status=?,updated_at=NOW() WHERE id=?`, m.ScheduledAt, m.HomeTeamID, m.AwayTeamID, m.HomeScore, m.AwayScore, m.Status, m.ID)
	return err
}

func (r *matchRepo) GetByID(id int64) (*models.Match, error) {
	m := &models.Match{}
	var deleted sql.NullTime
	row := r.db.QueryRow(`SELECT id,scheduled_at,home_team_id,away_team_id,home_score,away_score,status,created_at,updated_at,deleted_at FROM matches WHERE id=?`, id)
	if err := row.Scan(&m.ID, &m.ScheduledAt, &m.HomeTeamID, &m.AwayTeamID, &m.HomeScore, &m.AwayScore, &m.Status, &m.CreatedAt, &m.UpdatedAt, &deleted); err != nil {
		return nil, err
	}
	if deleted.Valid {
		t := deleted.Time
		m.DeletedAt = &t
	}
	return m, nil
}

func (r *matchRepo) ListByTeam(teamID int64) ([]*models.Match, error) {
	rows, err := r.db.Query(`SELECT id,scheduled_at,home_team_id,away_team_id,home_score,away_score,status,created_at,updated_at FROM matches WHERE (home_team_id=? OR away_team_id=?) AND deleted_at IS NULL ORDER BY scheduled_at DESC`, teamID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*models.Match
	for rows.Next() {
		m := &models.Match{}
		if err := rows.Scan(&m.ID, &m.ScheduledAt, &m.HomeTeamID, &m.AwayTeamID, &m.HomeScore, &m.AwayScore, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func (r *matchRepo) CountWins(teamID int64, asHome bool, upTo time.Time) (int, error) {
	var cnt int
	if asHome {

		err := r.db.QueryRow(`SELECT COUNT(1) FROM matches WHERE home_team_id=? AND home_score>away_score AND status='finished' AND scheduled_at<=?`, teamID, upTo).Scan(&cnt)
		return cnt, err
	}
	// as away
	err := r.db.QueryRow(`SELECT COUNT(1) FROM matches WHERE away_team_id=? AND away_score>home_score AND status='finished' AND scheduled_at<=?`, teamID, upTo).Scan(&cnt)
	return cnt, err
}

package repository

import (
	"database/sql"
	"time"
	"xyz_backend/src/models"
)

type GoalRepository interface {
	Create(g *models.Goal) error
	ListByMatch(matchID int64) ([]*models.Goal, error)
}

type goalRepo struct{ db *sql.DB }

func NewGoalRepo(db *sql.DB) GoalRepository { return &goalRepo{db} }

func (r *goalRepo) Create(g *models.Goal) error {
	res, err := r.db.Exec(`INSERT INTO goals (match_id,player_id,team_id,minute_scored) VALUES (?,?,?,?)`, g.MatchID, g.PlayerID, g.TeamID, g.MinuteScored)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	g.ID = id
	g.CreatedAt = time.Now()
	return nil
}

func (r *goalRepo) ListByMatch(matchID int64) ([]*models.Goal, error) {
	query := `	SELECT 
					g.id, g.match_id, g.player_id, p.team_id,
					g.created_at, g.minute_scored,
					p.name, p.jersey_number
				FROM goals g
				JOIN players p ON g.player_id = p.id
				WHERE g.match_id = ?
				ORDER BY g.minute_scored ASC
			`
	rows, err := r.db.Query(query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*models.Goal
	for rows.Next() {
		var g models.Goal
		err := rows.Scan(
			&g.ID,
			&g.MatchID,
			&g.PlayerID,
			&g.TeamID,
			&g.CreatedAt,
			&g.MinuteScored,
			&g.Player.Name,
			&g.Player.JerseyNumber,
		)
		if err != nil {
			return nil, err
		}
		goals = append(goals, &g)
	}
	return goals, nil
}

package usecase

import (
	"fmt"
	"xyz_backend/src/repository"
)

type ReportUsecase struct {
	matchRepo repository.MatchRepository
	goalRepo  repository.GoalRepository
	teamRepo  repository.TeamRepository
}

func NewReportUsecase(mr repository.MatchRepository, gr repository.GoalRepository, tr repository.TeamRepository) *ReportUsecase {
	return &ReportUsecase{matchRepo: mr, goalRepo: gr, teamRepo: tr}
}

func (u *ReportUsecase) MatchReport(matchID int64) (map[string]interface{}, error) {
	m, err := u.matchRepo.GetByID(matchID)
	if err != nil {
		return nil, err
	}
	goals, err := u.goalRepo.ListByMatch(matchID)
	if err != nil {
		return nil, err
	}

	scorerCount := map[int64]int{}
	for _, g := range goals {
		scorerCount[g.PlayerID]++
	}
	var topPlayerID int64
	var topGoals int
	for pid, cnt := range scorerCount {
		if cnt > topGoals {
			topGoals = cnt
			topPlayerID = pid
		}
	}

	status := "Draw"
	if m.HomeScore > m.AwayScore {
		status = "Tim Home Menang"
	}
	if m.AwayScore > m.HomeScore {
		status = "Tim Away Menang"
	}

	upTo := m.ScheduledAt
	homeWins, _ := u.matchRepo.CountWins(m.HomeTeamID, true, upTo)
	awayWins, _ := u.matchRepo.CountWins(m.AwayTeamID, false, upTo)

	return map[string]interface{}{
		"match":                m,
		"home_team":            m.HomeTeamID,
		"away_team":            m.AwayTeamID,
		"score":                fmt.Sprintf("%d-%d", m.HomeScore, m.AwayScore),
		"status_text":          status,
		"top_scorer_player_id": topPlayerID,
		"top_scorer_goals":     topGoals,
		"accum_home_wins":      homeWins,
		"accum_away_wins":      awayWins,
	}, nil
}

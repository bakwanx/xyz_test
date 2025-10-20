package usecase

import (
	"errors"
	"time"
	"xyz_backend/src/models"
	"xyz_backend/src/repository"
)

type MatchUsecase struct {
	matchRepo repository.MatchRepository
	teamRepo  repository.TeamRepository
	goalRepo  repository.GoalRepository
}

func NewMatchUsecase(mr repository.MatchRepository, tr repository.TeamRepository, gr repository.GoalRepository) *MatchUsecase {
	return &MatchUsecase{matchRepo: mr, teamRepo: tr, goalRepo: gr}
}

func (u *MatchUsecase) Create(m *models.Match) error {
	if _, err := u.teamRepo.GetByID(m.HomeTeamID); err != nil {
		return errors.New("home team not found")
	}
	if _, err := u.teamRepo.GetByID(m.AwayTeamID); err != nil {
		return errors.New("away team not found")
	}
	return u.matchRepo.Create(m)
}

func (u *MatchUsecase) ReportResult(matchID int64, homeScore, awayScore int, goals []*models.Goal) error {
	m, err := u.matchRepo.GetByID(matchID)
	if err != nil {
		return err
	}
	m.HomeScore = homeScore
	m.AwayScore = awayScore
	m.Status = "finished"
	if err := u.matchRepo.Update(m); err != nil {
		return err
	}
	for _, g := range goals {
		if err := u.goalRepo.Create(g); err != nil {
			return err
		}
	}
	return nil
}

func (u *MatchUsecase) GetMatchReport(matchID int64) (*models.Match, []*models.Goal, error) {
	m, err := u.matchRepo.GetByID(matchID)
	if err != nil {
		return nil, nil, err
	}
	goals, err := u.goalRepo.ListByMatch(matchID)
	if err != nil {
		return m, nil, err
	}
	return m, goals, nil
}

func (u *MatchUsecase) CountWins(teamID int64, asHome bool, upTo time.Time) (int, error) {
	return u.matchRepo.CountWins(teamID, asHome, upTo)
}

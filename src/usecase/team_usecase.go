package usecase

import (
	"errors"
	"xyz_backend/src/models"
	"xyz_backend/src/repository"
)

type TeamCreateInput struct {
	Name        string
	FoundedYear int
	Address     string
	City        string
	LogoPath    *string
}

type TeamUsecase struct {
	teamRepo repository.TeamRepository
}

func NewTeamUsecase(tr repository.TeamRepository) *TeamUsecase { return &TeamUsecase{teamRepo: tr} }

func (u *TeamUsecase) Create(in TeamCreateInput) (*models.Team, error) {
	if in.Name == "" {
		return nil, errors.New("name required")
	}
	t := &models.Team{
		Name:        in.Name,
		FoundedYear: in.FoundedYear,
		Address:     in.Address,
		City:        in.City,
	}
	if in.LogoPath != nil {
		t.LogoPath = in.LogoPath
	}
	if err := u.teamRepo.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (u *TeamUsecase) List(offset, limit int) ([]*models.Team, int, error) {
	teams, err := u.teamRepo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	total, err := u.teamRepo.Count()
	if err != nil {
		return nil, 0, err
	}
	return teams, total, nil
}

func (u *TeamUsecase) Delete(id int64) error {
	if checkTeam, err := u.teamRepo.GetByID(id); err != nil {
		if checkTeam == nil {
			return errors.New("Team not found")
		}
	}
	err := u.teamRepo.SoftDelete(id)

	return err
}

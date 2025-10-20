package usecase

import (
	"errors"
	"fmt"
	"xyz_backend/src/models"
	"xyz_backend/src/repository"
)

type PlayerUsecase struct {
	playerRepo repository.PlayerRepository
	teamRepo   repository.TeamRepository
}

func NewPlayerUsecase(pr repository.PlayerRepository, tr repository.TeamRepository) *PlayerUsecase {
	return &PlayerUsecase{playerRepo: pr, teamRepo: tr}
}

func (u *PlayerUsecase) Create(p *models.Player) error {
	if _, err := u.teamRepo.GetByID(p.TeamID); err != nil {
		return errors.New("team not found")
	}

	isAvailable, err := u.playerRepo.IsJerseyNumberAvailable(p.TeamID, p.JerseyNumber)
	if err != nil {
		return err
	}

	if !isAvailable {
		errMsg := fmt.Sprintf("Jersey with number %d has been used", +p.JerseyNumber)
		return errors.New(errMsg)
	}
	return u.playerRepo.Create(p)
}

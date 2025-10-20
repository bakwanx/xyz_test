package http_delivery

import (
	"net/http"
	"xyz_backend/src/models"
	"xyz_backend/src/usecase"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct{ playerUC *usecase.PlayerUsecase }

func NewPlayerHandler(uc *usecase.PlayerUsecase) *PlayerHandler { return &PlayerHandler{playerUC: uc} }

func (h *PlayerHandler) RegisterPlayer(c *gin.Context) {
	var req struct {
		TeamId       int64  `json:"team_id" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Height       int    `json:"height" binding:"required"`
		Weight       int    `json:"weight" binding:"required"`
		Position     string `json:"position" binding:"required"`
		JerseyNumber int    `json:"jersey_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.playerUC.Create(&models.Player{
		TeamID:       req.TeamId,
		Name:         req.Name,
		HeightCm:     req.Height,
		WeightKg:     req.Weight,
		Position:     req.Position,
		JerseyNumber: req.JerseyNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "player's added"})
}

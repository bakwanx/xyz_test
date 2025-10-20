package http_delivery

import (
	"net/http"
	"strconv"
	"time"
	"xyz_backend/src/models"
	"xyz_backend/src/usecase"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct{ matchUC *usecase.MatchUsecase }

func NewMatchHandler(uc *usecase.MatchUsecase) *MatchHandler { return &MatchHandler{matchUC: uc} }

func (h *MatchHandler) Create(c *gin.Context) {
	var req struct {
		ScheduledAt string `json:"scheduled_at" binding:"required"`
		HomeTeamID  int64  `json:"home_team_id" binding:"required"`
		AwayTeamID  int64  `json:"away_team_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	t, err := time.Parse(time.DateTime, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_at format, use DateTime"})
		return
	}

	m := &models.Match{ScheduledAt: t, HomeTeamID: req.HomeTeamID, AwayTeamID: req.AwayTeamID, Status: "scheduled"}
	if err := h.matchUC.Create(m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *MatchHandler) ReportResult(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req struct {
		HomeScore int `json:"home_score"`
		AwayScore int `json:"away_score"`
		Goals     []struct {
			PlayerID int64 `json:"player_id"`
			TeamID   int64 `json:"team_id"`
			Minute   int   `json:"minute"`
		} `json:"goals"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	var goals []*models.Goal
	for _, g := range req.Goals {
		goals = append(goals, &models.Goal{MatchID: id, PlayerID: g.PlayerID, TeamID: g.TeamID, MinuteScored: g.Minute})
	}
	if err := h.matchUC.ReportResult(id, req.HomeScore, req.AwayScore, goals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "result reported"})
}

func (h *MatchHandler) GetReport(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	m, goals, err := h.matchUC.GetMatchReport(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"match": m, "goals": goals})
}

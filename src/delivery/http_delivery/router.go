package http_delivery

import (
	"xyz_backend/config"
	"xyz_backend/src/delivery/middleware"
	"xyz_backend/src/repository"
	"xyz_backend/src/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	authUC *usecase.AuthUsecase,
	teamUC *usecase.TeamUsecase,
	playerUC *usecase.PlayerUsecase,
	matchUC *usecase.MatchUsecase,
	reportUC *usecase.ReportUsecase,
	tokenRepo repository.TokenBlacklistRepository,
	cfg *config.Config,
) {
	authHandler := NewAuthHandler(authUC)
	teamHandler := NewTeamHandler(teamUC, cfg.UploadDir)
	playerHandler := NewPlayerHandler(playerUC)
	matchHandler := NewMatchHandler(matchUC)
	reportHandler := NewReportHandler(reportUC)

	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)

	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth(cfg, tokenRepo))
	{
		protected.POST("/logout", authHandler.Logout)
		protected.POST("/teams", teamHandler.Create)
		protected.GET("/teams", teamHandler.List)
		protected.DELETE("/teams", teamHandler.Delete)
		protected.POST("/player", playerHandler.RegisterPlayer)
		protected.POST("/matches", matchHandler.Create)
		protected.POST("/matches/:id/result", matchHandler.ReportResult)
		protected.GET("/matches/:id/report", matchHandler.GetReport)
		protected.GET("/report/:id", reportHandler.GetReport)
	}
}

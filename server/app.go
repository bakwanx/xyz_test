package server

import (
	"database/sql"
	"log"
	"xyz_backend/config"
	"xyz_backend/src/delivery/http_delivery"
	"xyz_backend/src/repository"
	"xyz_backend/src/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func Run() {
	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(cfg.GinMode)

	// repos
	adminRepo := repository.NewAdminRepo(db)
	teamRepo := repository.NewTeamRepo(db)
	playerRepo := repository.NewPlayerRepo(db)
	matchRepo := repository.NewMatchRepo(db)
	goalRepo := repository.NewGoalRepo(db)
	tokenRepo := repository.NewTokenBlacklistRepo(db)

	// usecases
	authUC := usecase.NewAuthUsecase(adminRepo, tokenRepo, cfg)
	teamUC := usecase.NewTeamUsecase(teamRepo)
	playerUC := usecase.NewPlayerUsecase(playerRepo, teamRepo)
	matchUC := usecase.NewMatchUsecase(matchRepo, teamRepo, goalRepo)
	reportUC := usecase.NewReportUsecase(matchRepo, goalRepo, teamRepo)

	// handlers & router
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	// r.Use(gin.Logger())
	r.Use(gin.Recovery())

	http_delivery.RegisterRoutes(r, authUC, teamUC, playerUC, matchUC, reportUC, tokenRepo, cfg)

	log.Printf("listening on %s", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("server run: %v", err)
	}

}

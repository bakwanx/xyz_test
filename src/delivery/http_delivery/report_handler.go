package http_delivery

import (
	"net/http"
	"strconv"
	"xyz_backend/src/usecase"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct{ reportUC *usecase.ReportUsecase }

func NewReportHandler(uc *usecase.ReportUsecase) *ReportHandler { return &ReportHandler{reportUC: uc} }

func (h *ReportHandler) GetReport(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	result, err := h.reportUC.MatchReport(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": result})
}

package http_delivery

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"xyz_backend/src/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	teamUC    *usecase.TeamUsecase
	uploadDir string
}

func NewTeamHandler(uc *usecase.TeamUsecase, uploadDir string) *TeamHandler {
	os.MkdirAll(uploadDir, os.ModePerm)
	return &TeamHandler{teamUC: uc, uploadDir: uploadDir}
}

func (h *TeamHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	founded := c.PostForm("founded_year")
	address := c.PostForm("address")
	city := c.PostForm("city")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}

	var foundedYear int
	if founded != "" {
		if v, err := strconv.Atoi(founded); err == nil {
			foundedYear = v
		}
	}

	var logoPath *string
	file, err := c.FormFile("logo")
	if err == nil && file != nil {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		dst := filepath.Join(h.uploadDir, filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}
		tmp := dst
		logoPath = &tmp
	}

	inp := usecase.TeamCreateInput{
		Name:        name,
		FoundedYear: foundedYear,
		Address:     address,
		City:        city,
		LogoPath:    logoPath,
	}
	t, err := h.teamUC.Create(inp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": t})
}

func (h *TeamHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size

	teams, total, err := h.teamUC.List(offset, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": teams, "meta": gin.H{"page": page, "size": size, "total": total}})
}

func (h *TeamHandler) Delete(c *gin.Context) {

	id := c.DefaultQuery("id_team", "0")
	idTeam, _ := strconv.Atoi(id)

	err := h.teamUC.Delete(int64(idTeam))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

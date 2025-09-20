package handler

import (
	"net/http"
	"strings"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetMaterials(ctx *gin.Context) {
	var materials []ds.Material
	var err error
	creatorID := h.Repository.GetUser()

	materialSearch := ctx.Query("material_search")

	materialSearch = strings.TrimSpace(materialSearch)

	if materialSearch == "" {
		materials, err = h.Repository.GetMaterials()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	} else {
		materials, err = h.Repository.SearchMaterials(materialSearch)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	}
	currentExperiment, err := h.Repository.GetExperimentDraft(creatorID)
	experimentId := int(currentExperiment.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.HTML(http.StatusOK, "materials.html", gin.H{
		"materials":       materials,
		"count":           h.Repository.GetExperimentItemsCount(),
		"material_search": materialSearch,
		"experimentId":    experimentId,
	})
}
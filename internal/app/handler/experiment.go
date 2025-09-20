package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetExperiment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}
	experimentItems, experiment, err := h.Repository.GetExperiment(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "experiment.html", gin.H{
		"experimentItems": experimentItems,
		"experiment":      experiment,
		"count":           h.Repository.GetExperimentItemsCount(),
	})
}

func (h *Handler) SoftDeleteExperiment(ctx *gin.Context){
	idStr := ctx.Param("id")
	experimentId, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.SoftDeleteExperiment(experimentId)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/materials")
}
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteItemFromExperiment(ctx *gin.Context) {
	experimentID, err := strconv.Atoi(ctx.Param("experiment_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	materialID, err := strconv.Atoi(ctx.Param("material_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	experiment, err := h.Repository.DeleteItemFromExperiment(uint(experimentID), uint(materialID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(experiment)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, ds.ExperimentWithLogins{
		Experiment:     experiment,
		CreatorLogin:   creatorLogin,
		ModeratorLogin: moderatorLogin,
	})
}

type updateExperimentItemRequest struct {
	MaterialMass *float64 `json:"material_mass"`
	GasVolume    *float64 `json:"gas_volume"`
}

func (h *Handler) UpdateExperimentItem(ctx *gin.Context) {
	experimentID, err := strconv.Atoi(ctx.Param("experiment_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	materialID, err := strconv.Atoi(ctx.Param("material_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req updateExperimentItemRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updateData := ds.ExperimentItem{}
	if req.MaterialMass != nil {
		updateData.MaterialMass = *req.MaterialMass
	}
	if req.GasVolume != nil {
		updateData.GasVolume = *req.GasVolume
	}

	experimentItem, err := h.Repository.UpdateExperimentItem(
		uint(experimentID),
		uint(materialID),
		updateData,
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, experimentItem)
}
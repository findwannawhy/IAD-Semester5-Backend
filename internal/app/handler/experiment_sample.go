package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/findwannawhy/IAD-Semester5/internal/app/dto"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteSampleFromExperiment(ctx *gin.Context) {
  experimentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	sampleID, err := strconv.Atoi(ctx.Param("sample_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	uintSampleID := uint(sampleID)
	experiment, err := h.Repository.DeleteSampleFromExperiment(uint(experimentID), uintSampleID)
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

	creatorLogin, _, err := h.Repository.GetModeratorAndCreatorLogin(experiment)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.DeleteExperimentSampleResp{
		ExperimentID:   experiment.ID,
		SampleID:       uintSampleID,
		CreatorLogin:   creatorLogin,
	})
}

type updateExperimentSampleRequest struct {
	SampleMass       *float64 `json:"sample_mass"`
	EvolvedGasVolume *float64 `json:"evolved_gas_volume"`
}

func (h *Handler) UpdateExperimentSample(ctx *gin.Context) {
  experimentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	sampleID, err := strconv.Atoi(ctx.Param("sample_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req updateExperimentSampleRequest
	if err := ctx.BindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updateData := ds.ExperimentSample{}
	if req.SampleMass != nil {
		updateData.SampleMass = req.SampleMass
	}
	if req.EvolvedGasVolume != nil {
		updateData.EvolvedGasVolume = req.EvolvedGasVolume
	}

	experimentSample, err := h.Repository.UpdateExperimentSample(
		uint(experimentID),
		uint(sampleID),
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

	ctx.JSON(http.StatusOK, experimentSample)
}


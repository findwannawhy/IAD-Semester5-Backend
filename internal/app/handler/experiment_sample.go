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

// DeleteSampleFromExperiment godoc
// @Summary Удалить образец из исследования
// @Description Удаляет связь образца и исследования
// @Tags experiments-samples
// @Produce json
// @Param id path int true "ID исследования"
// @Param sample_id path int true "ID образца"
// @Success 200 {object} dto.ExperimentResponse "Обновленное исследование"
// @Failure 400 {object} map[string]string "Неверные ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments/{id}/soluble-samples/{sample_id} [delete]
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

	uintExperimentID := uint(experimentID)

	err = h.Repository.DeleteSampleFromExperiment(uintExperimentID, uint(sampleID))
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

	updatedExperiment, err := h.getExperimentData(uintExperimentID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedExperiment)
}

// UpdateExperimentSample godoc
// @Summary Изменить данные образца в исследовании
// @Description Обновляет параметры образца в конкретном исследовании
// @Tags experiments-samples
// @Accept json
// @Produce json
// @Param id path int true "ID исследования"
// @Param sample_id path int true "ID образца"
// @Param data body dto.UpdateExperimentSampleReq true "Новые данные"
// @Success 200 {object} ds.ExperimentSample "Обновленные данные"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments/{id}/soluble-samples/{sample_id} [put]
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

	var req dto.UpdateExperimentSampleReq
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


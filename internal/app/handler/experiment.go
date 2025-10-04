package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/findwannawhy/IAD-Semester5/internal/app/dto"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetExperiments(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from = time.Time{}
	var to = time.Time{}
	if fromDate != "" {
		from1, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = from1
	}
	fmt.Println(fromDate)

	toDate := ctx.Query("to-date")
	if toDate != "" {
		to1, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = to1
	}

	status := ctx.Query("status")

	experiments, err := h.Repository.GetExperiments(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := make([]dto.ExperimentsResponse, 0, len(experiments))
	for _, experiment := range experiments {
		creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(experiment)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		
		molarVolume := 0.0
		if experiment.MolarVolume != nil {
			molarVolume = *experiment.MolarVolume
		}
		
		resp = append(resp, dto.ExperimentsResponse{
			ID:             experiment.ID,
			MolarVolume:    molarVolume,
			Status:         experiment.Status,
			CreatedAt:      experiment.CreatedAt,
			FormedAt:       experiment.FormedAt,
			FinishedAt:     experiment.FinishedAt,
			CreatorLogin:   creatorLogin,
			ModeratorLogin: moderatorLogin,
		})
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetExperimentCart(ctx *gin.Context) {
	samplesCount := h.Repository.GetExperimentCount(h.Repository.GetUserID())

	if samplesCount == 0 {
		ctx.JSON(http.StatusOK, dto.DraftExperimentResponse{
			ExperimentID: 0,
			SampleCount:  0,
		})
		return
	}

	experiment, err := h.Repository.CheckCurrentExperimentDraft(h.Repository.GetUserID())
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, dto.DraftExperimentResponse{
				ExperimentID: 0,
				SampleCount:  0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.DraftExperimentResponse{
		ExperimentID: experiment.ID,
		SampleCount:  h.Repository.GetExperimentCount(experiment.CreatorID),
	})
}

func (h *Handler) GetExperiment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	samples, experiment, err := h.Repository.GetExperimentSamplesData(uint(id))
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

	// Получаем данные из experiments_samples для каждого образца
	experimentSamples, err := h.Repository.GetExperimentSamples(experiment.ID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Создаем мапу для быстрого доступа к данным ExperimentSample по SampleID
	experimentSampleMap := make(map[uint]ds.ExperimentSample)
	for _, es := range experimentSamples {
		experimentSampleMap[es.SampleID] = es
	}

	// Формируем ExperimentSampleCards
	experimentSampleCards := make([]dto.ExperimentSampleCard, 0, len(samples))
	for _, sample := range samples {
		es := experimentSampleMap[sample.ID]
		
		sampleMass := ""
		if es.SampleMass != nil {
			sampleMass = fmt.Sprintf("%.2f", *es.SampleMass)
		}
		
		evolvedGasVolume := ""
		if es.EvolvedGasVolume != nil {
			evolvedGasVolume = fmt.Sprintf("%.2f", *es.EvolvedGasVolume)
		}
		
		massFractionPercentage := ""
		if es.MassFractionPercentage != nil {
			massFractionPercentage = fmt.Sprintf("%.2f", *es.MassFractionPercentage)
		}

		imageURL := ""
		if sample.ImageURL != nil {
			imageURL = *sample.ImageURL
		}

		experimentSampleCards = append(experimentSampleCards, dto.ExperimentSampleCard{
			SampleID:                   sample.ID,
			Title:                      sample.Title,
			Formula:                    sample.Formula,
			ImageURL:                   imageURL,
			RelativeMolecularMass:      sample.RelativeMolecularMass,
			StoichiometricCoefficient:  sample.StoichiometricCoefficient,
			SampleMass:                 sampleMass,
			EvolvedGasVolume:           evolvedGasVolume,
			MassFractionPercentage:     massFractionPercentage,
		})
	}

	ctx.JSON(http.StatusOK, dto.ExperimentResponse{
		Experiment:            experiment,
		ExperimentSampleCards: experimentSampleCards,
		CreatorLogin:          creatorLogin,
		ModeratorLogin:        moderatorLogin,
	})
}

func (h *Handler) FormExperiment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "formed"

	experiment, err := h.Repository.FormExperiment(uint(id), status)
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

	ctx.JSON(http.StatusOK, dto.FormExperiment{
		Experiment:   experiment,
		CreatorLogin: creatorLogin,
	})
}

func (h *Handler) UpdateExperiment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var experimentJSON ds.ImpurityFractionExperiment
	if err := ctx.BindJSON(&experimentJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	experiment, err := h.Repository.UpdateExperiment(uint(id), experimentJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
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

	ctx.JSON(http.StatusOK, dto.UpdateExperiment{
		Experiment:   experiment,
		CreatorLogin: creatorLogin,
	})
}

func (h *Handler) SoftDeleteExperiment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	experimentId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "deleted"

	experiment, err := h.Repository.FormExperiment(uint(experimentId), status)
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

	ctx.JSON(http.StatusOK, dto.SoftDeleteExperiment{
		ExperimentID: experiment.ID,
		Status:       experiment.Status,
		FormedAt:     experiment.FormedAt,
		CreatorLogin: creatorLogin,
	})
}

func (h *Handler) ModerateExperiment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var statusJSON struct {
		Status string `json:"status"`
	}
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	experiment, err := h.Repository.ModerateExperiment(uint(id), statusJSON.Status)
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

	ctx.JSON(http.StatusOK, dto.ModerateExperiment{
		Experiment:     experiment,
		CreatorLogin:   creatorLogin,
		ModeratorLogin: moderatorLogin,
	})
}
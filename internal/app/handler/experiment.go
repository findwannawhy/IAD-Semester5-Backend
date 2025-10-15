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
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// GetExperiments godoc
// @Summary Получить список исследований
// @Description Возвращает исследования с возможностью фильтрации по датам и статусу
// @Tags impurity-experiments
// @Produce json
// @Param from-date query string false "Начальная дата (YYYY-MM-DD)"
// @Param to-date query string false "Конечная дата (YYYY-MM-DD)"
// @Param status query string false "Статус исследования"
// @Success 200 {array} dto.ExperimentsResponse "Список исследований"
// @Failure 400 {object} map[string]string "Неверный формат даты"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments [get]
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

	experiments = h.filterExperimentsByAuth(experiments, ctx)

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

// GetExperimentDraft godoc
// @Summary Получить черновик исследования
// @Description Возвращает информацию о текущем черновике исследования пользователя
// @Tags impurity-experiments
// @Produce json
// @Success 200 {object} dto.DraftExperimentResponse "Данные черновика исследования"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments/draft [get]
func (h *Handler) GetExperimentDraft(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	samplesCount := h.Repository.GetExperimentCount(userID)

	if samplesCount == 0 {
		ctx.JSON(http.StatusOK, dto.DraftExperimentResponse{
			ExperimentID: 0,
			SampleCount:  0,
		})
		return
	}

	experiment, err := h.Repository.CheckCurrentExperimentDraft(userID)
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

// GetExperiment godoc
// @Summary Получить исследование по ID
// @Description Возвращает полную информацию об исследовании включая образцы
// @Tags impurity-experiments
// @Produce json
// @Param id path int true "ID исследования"
// @Success 200 {object} dto.ExperimentResponse "Данные исследования с образцами"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments/{id} [get]
func (h *Handler) GetExperiment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := h.getExperimentData(uint(id))
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

	ctx.JSON(http.StatusOK, response)
}

// FormExperiment godoc
// @Summary Сформировать исследование
// @Description Переводит исследование в статус "formed"
// @Tags impurity-experiments
// @Produce json
// @Param id path int true "ID исследования"
// @Success 200 {object} dto.FormExperiment "Сформированное исследование"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments/{id}/form [put]
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

// UpdateExperiment godoc
// @Summary Изменить исследование
// @Description Обновляет данные исследования
// @Tags impurity-experiments
// @Accept json
// @Produce json
// @Param id path int true "ID исследования"
// @Param experiment body ds.ImpurityFractionExperiment true "Новые данные исследования"
// @Success 200 {object} dto.UpdateExperiment "Обновленное исследование"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments/{id} [put]
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

// SoftDeleteExperiment godoc
// @Summary Удалить исследование
// @Description Выполняет логическое удаление исследования
// @Tags impurity-experiments
// @Produce json
// @Param id path int true "ID исследования"
// @Success 200 {object} dto.SoftDeleteExperiment "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments/{id} [delete]
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

// ModerateExperiment godoc
// @Summary Модерировать исследование
// @Description Изменяет статус исследования (только для модераторов)
// @Tags impurity-experiments
// @Accept json
// @Produce json
// @Param id path int true "ID исследования"
// @Param status body dto.StatusJSON true "Новый статус"
// @Success 200 {object} dto.ExperimentResponse "Результат модерации"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /impurity-experiments/{id}/moderation [put]
func (h *Handler) ModerateExperiment(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
	}

	var statusJSON dto.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
					h.errorHandler(ctx, http.StatusNotFound, err)
			} else {
					h.errorHandler(ctx, http.StatusInternalServerError, err)
			}
			return
	}
	
	if !user.IsModerator {
			h.errorHandler(ctx, http.StatusForbidden, errors.New("требуются права модератора"))
			return
	}

	experiment, err := h.Repository.ModerateExperiment(uint(id), statusJSON.Status, userID)
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

	response, err := h.getExperimentData(experiment.ID)
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

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) filterExperimentsByAuth(experiments []ds.ImpurityFractionExperiment, ctx *gin.Context) []ds.ImpurityFractionExperiment {
userID, err := getUserID(ctx)
if err != nil {
	return []ds.ImpurityFractionExperiment{}
}

user, err := h.Repository.GetUserByID(userID)
if err == repository.ErrNotFound {
	return []ds.ImpurityFractionExperiment{}
}
if err != nil {
	return []ds.ImpurityFractionExperiment{}
}

if user.IsModerator {
	return experiments
}

var userExperiments []ds.ImpurityFractionExperiment
	for _, experiment := range experiments {
			fmt.Println(experiment.ID)
			if experiment.CreatorID == userID {
					userExperiments = append(userExperiments, experiment)
			}
	}
	
	return userExperiments

}

func (h *Handler) hasAccessToExperiment(creatorID uuid.UUID, ctx *gin.Context) bool {
userID, err := getUserID(ctx)
if err != nil {
	return false
}

user, err := h.Repository.GetUserByID(userID)
if err == repository.ErrNotFound {
	return false
}
if err != nil {
	return false
}

return creatorID == userID || user.IsModerator
}

func (h *Handler) getExperimentData(id uint) (dto.ExperimentResponse, error) {
	samples, experiment, err := h.Repository.GetExperimentSamplesData(id)
	if err != nil {
		return dto.ExperimentResponse{}, err
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(experiment)
	if err != nil {
		return dto.ExperimentResponse{}, err
	}

	experimentSamples, err := h.Repository.GetExperimentSamples(experiment.ID)
	if err != nil {
		return dto.ExperimentResponse{}, err
	}

	experimentSampleMap := make(map[uint]ds.ExperimentSample)
	for _, es := range experimentSamples {
		experimentSampleMap[es.SampleID] = es
	}

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

	return dto.ExperimentResponse{
		Experiment:            experiment,
		ExperimentSampleCards: experimentSampleCards,
		CreatorLogin:          creatorLogin,
		ModeratorLogin:        moderatorLogin,
	}, nil
}
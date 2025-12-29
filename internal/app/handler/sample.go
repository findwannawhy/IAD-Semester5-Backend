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

// GetSamples godoc
// @Summary Получить список образцов
// @Description Возвращает все образцы или фильтрует по названию. Поддерживает пагинацию.
// @Tags soluble-samples
// @Produce json
// @Param search_sample query string false "Название образца для поиска"
// @Param recently_viewed query boolean false "Показать недавно просмотренные"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество на странице (по умолчанию 20, макс 100)"
// @Success 200 {object} repository.PaginatedSamples "Список образцов с пагинацией"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /soluble-samples [get]
func (h *Handler) GetSamples(ctx *gin.Context) {
	var samples []ds.AcidSolubleSample
	var err error

	recentlyViewed := ctx.Query("recently_viewed")
	if recentlyViewed == "true" {
		// Получаем ID гостевой сессии
		sessionID, exists := ctx.Get("guest_session_id")
		if !exists {
			ctx.JSON(http.StatusOK, []ds.AcidSolubleSample{})
			return
		}

		sessionIDStr, ok := sessionID.(string)
		if !ok {
			ctx.JSON(http.StatusOK, []ds.AcidSolubleSample{})
			return
		}

		// Получаем список ID просмотренных образцов
		viewedIDs, err := h.Repository.GetViewedSamples(ctx.Request.Context(), sessionIDStr)
		if err != nil {
			ctx.JSON(http.StatusOK, []ds.AcidSolubleSample{})
			return
		}

		if len(viewedIDs) == 0 {
			ctx.JSON(http.StatusOK, []ds.AcidSolubleSample{})
			return
		}

		// Получаем полные данные образцов
		samples, err = h.Repository.GetSamplesByIDs(viewedIDs)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, samples)
		return
	}

	// Проверяем параметры пагинации
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	searchQuery := ctx.Query("search_sample")

	// Если указаны параметры пагинации, используем пагинированный метод
	if pageStr != "" || limitStr != "" {
		page := 1
		limit := 20

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
				if limit > 100 {
					limit = 100 // Максимум 100 записей на страницу
				}
			}
		}

		paginatedSamples, err := h.Repository.GetSamplesPaginated(searchQuery, page, limit)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, paginatedSamples)
		return
	}

	// Без пагинации (обратная совместимость)
	if searchQuery == "" {
		samples, err = h.Repository.GetSamples()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		samples, err = h.Repository.GetSamplesByName(searchQuery)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	if samples == nil {
		samples = make([]ds.AcidSolubleSample, 0)
	}
	ctx.JSON(http.StatusOK, samples)
}

// GetSample godoc
// @Summary Получить образец по ID
// @Description Возвращает информацию о образце по её идентификатору
// @Tags soluble-samples
// @Produce json
// @Param id path int true "ID образца"
// @Success 200 {object} ds.AcidSolubleSample "Данные образца"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Образец не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /soluble-samples/{id} [get]
func (h *Handler) GetSample(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	sample, err := h.Repository.GetSample(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// Записываем просмотр образца в гостевую сессию
	if sessionID, exists := ctx.Get("guest_session_id"); exists {
		if sessionIDStr, ok := sessionID.(string); ok {
			// Добавляем просмотренный образец, игнорируя ошибки
			_ = h.Repository.AddViewedSample(ctx.Request.Context(), sessionIDStr, id, 20*time.Minute)
		}
	}

	ctx.JSON(http.StatusOK, sample)
}

// CreateSample godoc
// @Summary Создать новый образец
// @Description Создает новый образец и возвращает его данные
// @Tags soluble-samples
// @Accept json
// @Produce json
// @Param sample body ds.AcidSolubleSample true "Данные нового образца"
// @Success 201 {object} ds.AcidSolubleSample "Созданный образец"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /soluble-samples [post]
func (h *Handler) CreateSample(ctx *gin.Context) {
	var sampleJSON ds.AcidSolubleSample
	if err := ctx.BindJSON(&sampleJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	sample, err := h.Repository.CreateSample(sampleJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/samples/%v", sample.ID))
	ctx.JSON(http.StatusCreated, sample)
}

// DeleteSample godoc
// @Summary Удалить образец
// @Description Выполняет логическое удаление образца по ID
// @Tags soluble-samples
// @Produce json
// @Param id path int true "ID образца"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Образец не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /soluble-samples/{id} [delete]
func (h *Handler) SoftDeleteSample(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	err = h.Repository.SoftDeleteSample(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sample_id": id,
		"status": "deleted",
	})
}

// UpdateSample godoc
// @Summary Изменить данные образца
// @Description Обновляет информацию о образце по ID
// @Tags soluble-samples
// @Accept json
// @Produce json
// @Param id path int true "ID образца"
// @Param sample body ds.AcidSolubleSample true "Новые данные образца"
// @Success 200 {object} ds.AcidSolubleSample "Обновленный образец"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Образец не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /soluble-samples/{id} [put]
func (h *Handler) UpdateSample(ctx *gin.Context) {
	var sampleJSON ds.AcidSolubleSample
	if err := ctx.BindJSON(&sampleJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id := uint(id64)

	sample, err := h.Repository.UpdateSample(id, sampleJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, sample)
}

// AddSampleToExperimentDraft godoc
// @Summary Добавить образец в черновик исследования
// @Description Добавляет образец в черновик исследования пользователя
// @Tags soluble-samples
// @Produce json
// @Param id path int true "ID образца"
// @Success 200 {object} dto.ExperimentResponse "Исследование с добавленным образцом"
// @Success 201 {object} dto.ExperimentResponse "Создано новое исследование"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Образец не найден"
// @Failure 409 {object} map[string]string "Образец уже в исследовании"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /soluble-samples/{id}/experiments/draft [post]
func (h *Handler) AddSampleToExperimentDraft(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	experiment, created, err := h.Repository.GetExperimentDraft(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	experimentId := experiment.ID

	sampleId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddSampleToExperiment(uint(experimentId), uint(sampleId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	status := http.StatusOK

	if created {
		ctx.Header("Location", fmt.Sprintf("/experiment/%v", experiment.ID))
		status = http.StatusCreated
	}

	updatedExperiment, err := h.getExperimentData(experiment.ID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}	

	ctx.JSON(status, updatedExperiment)
}

// UpdateImage godoc
// @Summary Загрузить изображение для образца
// @Description Загружает изображение для образца и возвращает обновленные данные
// @Tags soluble-samples
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID образца"
// @Param image formData file true "Изображение образца"
// @Success 200 {object} map[string]interface{} "Статус загрузки и данные образца"
// @Failure 400 {object} map[string]string "Неверный запрос или файл"
// @Failure 404 {object} map[string]string "Образец не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /soluble-samples/{id}/image [post]
func (h *Handler) UpdateImage(ctx *gin.Context) {
	sampleId64, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	sampleId := uint(sampleId64)

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	sample, err := h.Repository.UpdateImage(ctx, sampleId, file)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.UploadImage{
		SampleID: sample.ID,
		SampleTitle: sample.Title,
		ImageURL: *sample.ImageURL,
	})
}




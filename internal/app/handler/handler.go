package handler

import (
	"errors"

	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/api/acid-soluble-samples", h.GetSamples) // список образцов
	router.GET("/api/acid-soluble-samples/:id", h.GetSample) // получить образец по id
	router.POST("/api/acid-soluble-samples", h.CreateSample) // создать образец
	router.DELETE("/api/acid-soluble-samples/:id", h.SoftDeleteSample) // удалить образец (soft)
	router.PUT("/api/acid-soluble-samples/:id", h.UpdateSample) // обновить образец
	router.POST("/api/acid-soluble-samples/:id/experiments/draft", h.AddSampleToExperimentDraft) // добавить образец в корзину
	router.POST("/api/acid-soluble-samples/:id/image", h.UpdateImage) // обновить изображение образца

	router.GET("/api/impurity-fraction-experiments/draft", h.GetExperimentCart)	// текущий черновик пользователя
	router.GET("/api/impurity-fraction-experiments", h.GetExperiments) // список экспериментов
	router.GET("/api/impurity-fraction-experiments/:id", h.GetExperiment) // получить эксперимент по id
	router.PUT("/api/impurity-fraction-experiments/:id", h.UpdateExperiment) // обновить эксперимент
	router.PUT("/api/impurity-fraction-experiments/:id/formation", h.FormExperiment) // поменять статус эксперимента (user)
	router.PUT("/api/impurity-fraction-experiments/:id/moderation", h.ModerateExperiment) // поменять статус эксперимента (moderator)
	router.DELETE("/api/impurity-fraction-experiments/:id", h.SoftDeleteExperiment) // удалить эксперимент (soft)

	router.DELETE("/api/impurity-fraction-experiments/:id/samples/:sample_id", h.DeleteSampleFromExperiment) // удалить образец из эксперимента
	router.PUT("/api/impurity-fraction-experiments/:id/samples/:sample_id", h.UpdateExperimentSample) // обновить позицию образца в эксперименте

	router.POST("/api/users", h.CreateUser) // регистрация пользователя
	router.GET("/api/users/me", h.GetProfile) // профиль текущего пользователя
	router.PUT("/api/users/me", h.UpdateProfile) // обновить профиль
	router.POST("/api/users/login", h.SignIn) // вход в систему
	router.POST("/api/users/logout", h.SignOut) // выход из системы
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	
	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}
	
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}
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
	router.GET("/api/materials", h.GetMaterials) // список материалов
	router.GET("/api/materials/:id", h.GetMaterial) // получить материал по id
	router.POST("/api/materials", h.CreateMaterial) // создать материал
	router.DELETE("/api/materials/:id", h.SoftDeleteMaterial) // удалить материал (soft)
	router.PUT("/api/materials/:id", h.UpdateMaterial) // обновить материал
	router.POST("/api/materials/:id/draft", h.AddMaterialToExperiment) // добавить материал в корзину
	router.POST("/api/materials/:id/image", h.UpdateImage) // обновить изображение материала

	router.GET("/api/experiments/draft", h.GetExperimentCart)	// текущий черновик пользователя
	router.GET("/api/experiments", h.GetExperiments) // список экспериментов
	router.GET("/api/experiments/:id", h.GetExperiment) // получить эксперимент по id
	router.PUT("/api/experiments/:id", h.UpdateExperiment) // обновить эксперимент
	router.PUT("/api/experiments/:id/form", h.FormExperiment) // поменять статус эксперимента (user)
	router.PUT("/api/experiments/:id/moderation", h.ModerateExperiment) // поменять статус эксперимента (moderator)
	router.DELETE("/api/experiments/:id", h.SoftDeleteExperiment) // удалить эксперимент (soft)

	router.DELETE("/api/experiments/:id/materials/:material_id", h.DeleteItemFromExperiment) // удалить материал из эксперимента
	router.PUT("/api/experiments/:id/materials/:material_id", h.UpdateExperimentItem) // обновить позицию материала в эксперименте

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
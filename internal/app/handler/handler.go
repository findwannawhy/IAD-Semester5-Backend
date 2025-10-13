package handler

import (
	"errors"
	"net/http"

	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler godoc
// @title IAD API
// @version 1.0
// @description API для управления экспериментами по определению массовой доли примесей в образце
// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@impurity.fraction.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api/v1")

	unauthorized := api.Group("/")

	unauthorized.POST("/users/sign-up", h.SignUp)
	unauthorized.POST("/users/sign-in", h.SignIn)
	unauthorized.GET("/soluble-samples", h.GetSamples) // список образцов
	unauthorized.GET("/soluble-samples/:id", h.GetSample) // получить образец по id

	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))

	authorized.POST("/soluble-samples", h.CreateSample) // создать образец
	authorized.PUT("/soluble-samples/:id", h.UpdateSample) // обновить образец
	authorized.DELETE("/soluble-samples/:id", h.SoftDeleteSample) // удалить образец (soft)
	authorized.POST("/soluble-samples/:id/impurity-experiments/draft", h.AddSampleToExperimentDraft) // добавить образец в корзину
	authorized.POST("/soluble-samples/:id/image", h.UpdateImage) // обновить изображение образца

	authorized.GET("/impurity-experiments/draft", h.GetExperimentDraft)	// текущий черновик пользователя
	authorized.GET("/impurity-experiments", h.GetExperiments) // список экспериментов
	authorized.GET("/impurity-experiments/:id", h.GetExperiment) // получить эксперимент по id
	authorized.PUT("/impurity-experiments/:id", h.UpdateExperiment) // обновить эксперимент
	authorized.PUT("/impurity-experiments/:id/formation", h.FormExperiment) // поменять статус эксперимента (user)
	authorized.DELETE("/impurity-experiments/:id", h.SoftDeleteExperiment) // удалить эксперимент (soft)

	authorized.DELETE("/impurity-experiments/:id/soluble-samples/:sample_id", h.DeleteSampleFromExperiment) // удалить образец из эксперимента
	authorized.PUT("/impurity-experiments/:id/soluble-samples/:sample_id", h.UpdateExperimentSample) // обновить позицию образца в эксперименте

	authorized.GET("/users/me", h.GetProfile) // профиль текущего пользователя
	authorized.PUT("/users/me", h.UpdateProfile) // обновить профиль
	authorized.POST("/users/sign-out", h.SignOut) // выход из системы

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	moderator.PUT("/impurity-experiments/:id/moderation", h.ModerateExperiment) // поменять статус эксперимента (moderator)

	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
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
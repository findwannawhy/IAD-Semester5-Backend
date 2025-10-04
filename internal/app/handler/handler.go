package handler

import (
	"html/template"

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
	router.GET("/acid-soluble-samples", h.GetSamples)
	router.GET("/acid-soluble-samples/:id", h.GetSample)
	router.POST("/acid-soluble-samples/:id/experiment", h.AddSampleToExperiment)
	router.GET("/impurity-fraction-experiments/:id", h.GetExperiment)
	router.POST("/impurity-fraction-experiments/:id/delete", h.SoftDeleteExperiment)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	// Добавляем пользовательскую функцию для разыменования указателей
	router.SetFuncMap(template.FuncMap{
		"deref": func(ptr *float64) float64 {
			if ptr == nil {
				return 0
			}
			return *ptr
		},
	})
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
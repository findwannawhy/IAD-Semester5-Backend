package handler

import (
	"net/http"
	"strconv"

	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Обработчики

type Handler struct {
  Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
  return &Handler{
    Repository: r,
  }
}

// Материалы

func (h *Handler) GetMaterials(ctx *gin.Context) {
	var materials []repository.Material
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		materials, err = h.Repository.GetMaterials()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		materials, err = h.Repository.GetMaterialsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "materials.html", gin.H{
		"materials": materials,
		"query":  searchQuery,
	})
}

// Материал

func (h *Handler) GetMaterial(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id материала из урла
	id, err := strconv.Atoi(idStr) // форматируем в int
	if err != nil {
		logrus.Error(err)
	}

	material, err := h.Repository.GetMaterial(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "material.html", gin.H{
		"material": material,
	})
}

// Заявка

func (h *Handler) GetOrder(ctx *gin.Context) {
	order, err := h.Repository.GetOrder()
	if err != nil {
		logrus.Error(err)
	}

	orderItems, err := h.Repository.GetOrderItemsByOrderId(order.ID)
	if err != nil {
		logrus.Error(err)
	}

	acid, err := h.Repository.GetAcid(order.AcidId)
	if err != nil {
		logrus.Error(err)
	}

	materials, err := h.Repository.GetMaterials()
	if err != nil {
		logrus.Error(err)
	}
	materialsMap := make(map[int]repository.Material, len(materials))
	for _, m := range materials {
		materialsMap[m.ID] = m
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order": order,
		"orderItems": orderItems,
		"acid": acid,
		"materialsMap": materialsMap,
	})
}
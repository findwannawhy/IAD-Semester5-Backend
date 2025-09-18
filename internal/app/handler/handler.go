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

// View-модель для карточек эксперимента
type ExperimentCard struct {
  MaterialID int
  Title string
  Description string
  RelativeMolecularMass float64
  ImageURL string
  MaterialMass float64
  GasVolume float64
  MassFractionPercentage float64
  StoichiometricCoefficient float64
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

	// Обрабатываем параметр material_search из url, выводим пользователю только те материалы,
	// которые содержат в названии то, что ввёл пользователь в поиске
	materialSearch := ctx.Query("material_search")
	if materialSearch == "" {
		materials, err = h.Repository.GetMaterials()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		materials, err = h.Repository.GetMaterialsByTitle(materialSearch)
		if err != nil {
			logrus.Error(err)
		}
	}

	// В лаб1 у нас только один эксперимент, поэтому считаем все услуги всех экспериментов,
	// чтобы не добавлять id в url, ведущую на каталог. Потом, когда мы сможем получать
	// id эксперимента, на который будет вести корзина, мы заменим используемый метод на 
	// GetExperimentItemsCount(id int), где id - id эксперимента, на который будет вести корзина
	count, err := h.Repository.GetAllExperimentItemsCount()
	if err != nil {
		logrus.Error(err)
		count = 0
	}

  // Возвращаем html шаблон, передаём в него все материалы и параметр material_search,
	// чтобы можно было сохранить введённый поиск
	ctx.HTML(http.StatusOK, "materials.html", gin.H{
		"materials": materials,
		"material_search":  materialSearch,
		"Count": count,
	})
}

// Материал

func (h *Handler) GetMaterial(ctx *gin.Context) {
	// Получаем id материала из url
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	// Получаем материал из репозитория по id
	material, err := h.Repository.GetMaterial(id)
	if err != nil {
		logrus.Error(err)
	}

	// Возвращаем html шаблон
	ctx.HTML(http.StatusOK, "material.html", gin.H{
		"material": material,
	})
}

// Эксперимент

func (h *Handler) GetExperiment(ctx *gin.Context) {
	// Получаем id эксперимента из url
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	// Получаем эксперимент из репозитория по id
	experiment, err := h.Repository.GetExperiment(id)
	if err != nil {
		logrus.Error(err)
	}

	// Получаем список веществ, добавленных в эксперимент
	experimentItems := experiment.ExperimentItems

	// Получаем количество веществ, добавленных в эксперимент
	itemsAmount, err := h.Repository.GetExperimentItemsCount(id)
	if err != nil {
		logrus.Error(err)
	}

	// Получаем список материалов
	materials, err := h.Repository.GetMaterials()
	if err != nil {
		logrus.Error(err)
	}

	// Создаем map, чтобы можно было получить материал по его id
	materialsMap := make(map[int]repository.Material, len(materials))
	for _, m := range materials {
		materialsMap[m.ID] = m
	}

	// Собираем карточки для шаблона
	cards := make([]ExperimentCard, 0, len(experimentItems))
	for _, it := range experimentItems {
		m := materialsMap[it.MaterialId]
		cards = append(cards, ExperimentCard{
			MaterialID:                it.MaterialId,
			Title:                     m.Title,
			Description:               m.Description,
			RelativeMolecularMass:     m.RelativeMolecularMass,
			ImageURL:                  m.ImageURL,
			MaterialMass:              it.MaterialMass,
			GasVolume:                 it.GasVolume,
			MassFractionPercentage:    it.MassFractionPercentage,
			StoichiometricCoefficient: m.StoichiometricCoefficient,
		})
	}

	// Возвращаем html шаблон
	ctx.HTML(http.StatusOK, "experiment.html", gin.H{
		"experiment": experiment,
		"cards": cards,
		"itemsAmount": itemsAmount,
	})
}
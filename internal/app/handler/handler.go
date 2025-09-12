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

// View-модель для карточек заказа
type OrderCard struct {
  MaterialID int
  Title string
  Description string
  RelativeMolecularMass float64
  ImageURL string
  MaterialMass float64
  GasVolume float64
  MassFractionPercentage float64
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

	// Обрабатываем параметр query из url, выводим пользователю только те материалы,
	// которые содержат в названии то, что ввёл пользователь в поиске
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

	// В лаб1 у нас только один заказ, поэтому считаем все услуги всех заказов,
	// чтобы не добавлять id в url, ведущую на каталог. Потом, когда мы сможем получать
	// id заказа, на который будет вести корзина, мы заменим используемый метод
	count, err := h.Repository.GetAllOrderItemsCount()
	if err != nil {
		logrus.Error(err)
		count = 0
	}

  // Возвращаем html шаблон, передаём в него все материалы и параметр query,
	// чтобы можно было сохранить введённый поиск
	ctx.HTML(http.StatusOK, "materials.html", gin.H{
		"materials": materials,
		"query":  searchQuery,
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

// Заявка

func (h *Handler) GetOrder(ctx *gin.Context) {
	// Получаем id заказа из url
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	// Получаем заказ из репозитория по id
	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	// Получаем карточки, добавленные в заказ по id заказа
	orderItems, err := h.Repository.GetOrderItems(order.ID)
	if err != nil {
		logrus.Error(err)
	}

	// Получаем кислоты из репозитория, чтобы отобразить список для выбора в форме
	acids, err := h.Repository.GetAcids()
	if err != nil {
		logrus.Error(err)
	}

	// Получаем кислоту, которую выбрал пользователь в заказе, из уже загруженного списка
	var acid repository.Acid
	for _, a := range acids {
		if a.ID == order.AcidId {
			acid = a
			break
		}
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
	cards := make([]OrderCard, 0, len(orderItems))
	for _, it := range orderItems {
		m := materialsMap[it.MaterialId]
		cards = append(cards, OrderCard{
			MaterialID: it.MaterialId,
			Title: m.Title,
			Description: m.Description,
			RelativeMolecularMass: m.RelativeMolecularMass,
			ImageURL: m.ImageURL,
			MaterialMass: it.MaterialMass,
			GasVolume: it.GasVolume,
			MassFractionPercentage: it.MassFractionPercentage,
		})
	}

	// Возвращаем html шаблон
	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order": order,
		"acid": acid,
		"cards": cards,
	})
}
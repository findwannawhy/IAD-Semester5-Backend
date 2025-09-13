package repository

import (
	"fmt"
	"strings"
)

// Репозиторий

type Repository struct {
}

func NewRepository() (*Repository, error) {
  return &Repository{}, nil
}

// Материалы

type Material struct { 
  ID    int
  Title string
	Formula string
	Description string
	RelativeMolecularMass float64
	ImageURL string
	StoichiometricCoefficient float64
}

// Получаем все материалы
func (r *Repository) GetMaterials() ([]Material, error) {
  materials := []Material{
    {
      ID:    1,
      Title: "Известняк",
			Formula: "CaCO3",
			Description: "Осадочная порода, состоящая преимущественно из кальцита (карбоната кальция).",
			RelativeMolecularMass:  100.07,
			StoichiometricCoefficient: 1,
			ImageURL: "http://localhost:9000/img/izvestnyak.jpg",
    },
    {
      ID:    2,
      Title: "Мрамор",
			Formula: "CaCO3",
			Description: "Метаморфическая порода из кальцита, прочная, декоративная, полируемая.",
			RelativeMolecularMass:  100.07,
			StoichiometricCoefficient: 1,
			ImageURL: "http://localhost:9000/img/mramor.jpg",
    },
    {
      ID:    3,
      Title: "Металлический цинк",
			Formula: "Zn",
			Description: "Голубовато-белый металл, пластичный, коррозионно-стойкий.",
			RelativeMolecularMass:  65.39,
			StoichiometricCoefficient: 1,
			ImageURL: "http://localhost:9000/img/cink.jpg",
    },
		{
      ID:    4,
      Title: "Сода",
			Formula: "Na2CO3",
			Description: "Белый, без запаха, водорастворимый порошок или кристаллы, представляющие собой среднюю соль угольной кислоты",
			RelativeMolecularMass:  105.99,
			StoichiometricCoefficient: 1,
			ImageURL: "http://localhost:9000/img/soda.jpg",
    },
  }
  
  if len(materials) == 0 {
    return nil, fmt.Errorf("материалы не найдены")
  }

  return materials, nil
}

// Получаем материал по его id
func (r *Repository) GetMaterial(id int) (Material, error) {
	materials, err := r.GetMaterials()
	if err != nil {
		return Material{}, err
	}

	for _, material := range materials {
		if material.ID == id {
			return material, nil
		}
	}
	return Material{}, fmt.Errorf("материал не найден")
}

// Получаем материалы по названию
func (r *Repository) GetMaterialsByTitle(title string) ([]Material, error) {
	materials, err := r.GetMaterials()
	if err != nil {
		return []Material{}, err
	}

	// Убираем пробелы и нормализуем регистр для поискового запроса
	trimmedQuery := strings.TrimSpace(strings.ToLower(title))

	var result []Material
	for _, material := range materials {
		if strings.Contains(strings.ToLower(material.Title), trimmedQuery) {
			result = append(result, material)
		} else if strings.Contains(strings.ToLower(material.Formula), trimmedQuery) {
			result = append(result, material)
		}
	}

	return result, nil
}

// Заказ

type Order struct {
	ID int
	MolarVolume float64
}

// Получаем все заказы
func (r *Repository) GetOrders() ([]Order, error) {
  orders := []Order{
		{
			ID: 1,
			MolarVolume: 22.4,
		},
	}
  
  if len(orders) == 0 {
    return nil, fmt.Errorf("заказы не найдены")
  }

  return orders, nil
}

// Получаем заказ по его id
func (r *Repository) GetOrder(id int) (Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err
	}
	
	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("заказ не найден")
}

// Карточки заказов

type OrderItem struct {
	ID int
	OrderId int
	MaterialId int
	MaterialMass float64
	GasVolume float64
	MassFractionPercentage float64
}

// Получаем все услуги добавленные во все заказы
func (r *Repository) GetAllOrderItems() ([]OrderItem, error) {
  orderItems := []OrderItem{
    {
      ID:    1,
			OrderId: 1,
      MaterialId: 1,
      MaterialMass: 2.35,
      GasVolume: 0.484,
      MassFractionPercentage: 8,
    },
		{
      ID:    2,
			OrderId: 1,
      MaterialId: 2,
      MaterialMass: 1.8,
      GasVolume: 0.397,
      MassFractionPercentage: 1.5,
    },
  }
  
  if len(orderItems) == 0 {
    return nil, fmt.Errorf("не найдено заказанных услуг")
  }

  return orderItems, nil
}

// Получаем услуги добавленные в один заказ по id заказа
func (r *Repository) GetOrderItems(id int) ([]OrderItem, error) {
	orderItems, err := r.GetAllOrderItems()
	if err != nil {
		return []OrderItem{}, err
	}

	var result []OrderItem
	for _, orderItem := range orderItems {
		if orderItem.OrderId == id {
			result = append(result, orderItem)	
		}
	}
	return result, nil
}

// Получаем количество услуг добавленных во все заказы
func (r *Repository) GetAllOrderItemsCount() (int, error) {
	orderItems, err := r.GetAllOrderItems()
	if err != nil {
		return 0, err
	}

	return len(orderItems), nil
}

// Получаем количество услуг добавленных в один заказ по id заказа
func (r *Repository) GetOrderItemsCount(orderID int) (int, error) {
	orderItems, err := r.GetOrderItems(orderID)
	if err != nil {
		return 0, err
	}

	return len(orderItems), nil
}
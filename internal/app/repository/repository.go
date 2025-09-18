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

// Эксперимент

type ExperimentItem struct {
	ID int
	ExperimentId int
	MaterialId int
	MaterialMass float64
	GasVolume float64
	MassFractionPercentage float64
}

type Experiment struct {
	ID int
	MolarVolume float64
	ExperimentItems []ExperimentItem
}

// Получаем все эксперименты
func (r *Repository) GetAllExperiments() ([]Experiment, error) {
  experiments := []Experiment{
		{
			ID: 1,
			MolarVolume: 22.4,
			ExperimentItems: []ExperimentItem{
				{
					ID:    1,
					ExperimentId: 1,
					MaterialId: 1,
					MaterialMass: 2.35,
					GasVolume: 0.484,
					MassFractionPercentage: 8,
				},
				{
					ID:    2,
					ExperimentId: 1,
					MaterialId: 2,
					MaterialMass: 1.8,
					GasVolume: 0.397,
					MassFractionPercentage: 1.5,
				},
			},
		},
	}
  
  if len(experiments) == 0 {
    return nil, fmt.Errorf("эксперименты не найдены")
  }

  return experiments, nil
}

// Получаем эксперимент по его id
func (r *Repository) GetExperiment(id int) (Experiment, error) {
	experiments, err := r.GetAllExperiments()
	if err != nil {
		return Experiment{}, err
	}
	
	for _, experiment := range experiments {
		if experiment.ID == id {
			return experiment, nil
		}
	}
	return Experiment{}, fmt.Errorf("эксперимент не найден")
}

// Получаем количество услуг добавленных во все эксперименты
func (r *Repository) GetAllExperimentItemsCount() (int, error) {
	experiments, err := r.GetAllExperiments()
	if err != nil {
		return 0, err
	}

	totalItems := 0
	for _, experiment := range experiments {
		totalItems += len(experiment.ExperimentItems)
	}

	return totalItems, nil
}

// Получаем количество услуг добавленных в один эксперимент по id эксперимента
func (r *Repository) GetExperimentItemsCount(experimentID int) (int, error) {
	experiment, err := r.GetExperiment(experimentID)
	if err != nil {
		return 0, err
	}

	return len(experiment.ExperimentItems), nil
}
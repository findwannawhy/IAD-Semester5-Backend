package repository

import (
	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
)

func (r *Repository) GetMaterials() ([]ds.Material, error) {
	var materials []ds.Material
	err := r.db.Order("id").Where("deleted = false").Find(&materials).Error
	if err != nil {
		return nil, err
	}

	return materials, nil
}

func (r *Repository) GetMaterial(id int) (*ds.Material, error) {
	material := ds.Material{}
	err := r.db.Order("id").Where("id = ? and deleted = ?", id, false).First(&material).Error
	if err != nil {
		return &ds.Material{}, err
	}
	return &material, nil
}

func (r *Repository) SearchMaterials(materialSearch string) ([]ds.Material, error) {
	var materials []ds.Material
	err := r.db.Order("id").Where("(title ILIKE ? OR formula ILIKE ?) AND deleted = ?", "%"+materialSearch+"%", "%"+materialSearch+"%", false).Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}

func (r *Repository) AddMaterialToExperiment(experimentId int, materialId int) error {
	var material ds.Material
	if err := r.db.First(&material, materialId).Error; err != nil {
		return err
	}

	var experiment ds.Experiment
	if err := r.db.First(&experiment, experimentId).Error; err != nil {
		return err
	}
	
	experimentItems := ds.ExperimentItem{}
	result := r.db.Where("experiment_id = ? and material_id = ?", experimentId, materialId).Find(&experimentItems)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return nil
	}
	return r.db.Create(&ds.ExperimentItem{
		ExperimentID: experimentId,
		MaterialID: materialId,
	}).Error
}

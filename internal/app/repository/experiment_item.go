package repository

import (
	"errors"
	"fmt"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) DeleteItemFromExperiment(experimentId uint, materialId uint) (ds.Experiment, error) {
	var dbExperiment ds.Experiment
	err := r.db.Where("id = ?", experimentId).First(&dbExperiment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Experiment{}, fmt.Errorf("%w: исследование с id %d", ErrNotFound, experimentId)
		}
		return ds.Experiment{}, err
	}

	err = r.db.Where("material_id = ? and experiment_id = ?", materialId, experimentId).Delete(&ds.ExperimentItem{}).Error
	if err != nil {
		return ds.Experiment{}, err
	}
	return dbExperiment, nil
}

func (r *Repository) UpdateExperimentItem(experimentId uint, materialId uint, ExperimentItem ds.ExperimentItem) (ds.ExperimentItem, error) {
	var experimentItem ds.ExperimentItem
	err := r.db.Model(&experimentItem).Where("material_id = ? and experiment_id = ?", materialId, experimentId).Updates(experimentItem).First(&experimentItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ExperimentItem{}, fmt.Errorf("%w: материал в эксперименте", ErrNotFound)
		}
		return ds.ExperimentItem{}, err
	}
	return experimentItem, nil
}
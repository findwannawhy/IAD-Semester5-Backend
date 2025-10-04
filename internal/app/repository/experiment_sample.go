package repository

import (
	"errors"
	"fmt"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) DeleteSampleFromExperiment(experimentId uint, sampleId uint) (ds.ImpurityFractionExperiment, error) {
	var dbExperiment ds.ImpurityFractionExperiment
	err := r.db.Where("id = ?", experimentId).First(&dbExperiment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ImpurityFractionExperiment{}, fmt.Errorf("%w: эксперимент с id %d", ErrNotFound, experimentId)
		}
		return ds.ImpurityFractionExperiment{}, err
	}

	result := r.db.Where("sample_id = ? and experiment_id = ?", sampleId, experimentId).Delete(&ds.ExperimentSample{})
	if result.Error != nil {
		return ds.ImpurityFractionExperiment{}, result.Error
	}
	if result.RowsAffected == 0 {
		return ds.ImpurityFractionExperiment{}, fmt.Errorf("%w: образец с id %d не найден в эксперименте с id %d", ErrNotFound, sampleId, experimentId)
	}
	return dbExperiment, nil
}

func (r *Repository) UpdateExperimentSample(experimentId uint, sampleId uint, experimentSample ds.ExperimentSample) (ds.ExperimentSample, error) {
	var expSample ds.ExperimentSample
	err := r.db.Model(&expSample).Where("sample_id = ? and experiment_id = ?", sampleId, experimentId).Updates(experimentSample).First(&expSample).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ExperimentSample{}, fmt.Errorf("%w: образец в эксперименте", ErrNotFound)
		}
		return ds.ExperimentSample{}, err
	}
	return expSample, nil
}


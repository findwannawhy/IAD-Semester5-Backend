package repository

import (
	"errors"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"gorm.io/gorm"
)

func (r *Repository) GetSamples() ([]ds.AcidSolubleSample, error) {
	var samples []ds.AcidSolubleSample
	err := r.db.Order("id").Where("deleted = false").Find(&samples).Error
	if err != nil {
		return nil, err
	}

	return samples, nil
}

func (r *Repository) GetSample(id int) (*ds.AcidSolubleSample, error) {
	sample := ds.AcidSolubleSample{}
	err := r.db.Order("id").Where("id = ? and deleted = ?", id, false).First(&sample).Error
	if err != nil {
		return &ds.AcidSolubleSample{}, err
	}
	return &sample, nil
}

func (r *Repository) SearchSamples(searchSample string) ([]ds.AcidSolubleSample, error) {
	var samples []ds.AcidSolubleSample
	err := r.db.Order("id").Where("(title ILIKE ? OR formula ILIKE ?) AND deleted = ?", "%"+searchSample+"%", "%"+searchSample+"%", false).Find(&samples).Error
	if err != nil {
		return nil, err
	}
	return samples, nil
}

func (r *Repository) AddSampleToExperiment(experimentId uint, sampleId uint) error {
	var sample ds.AcidSolubleSample
	if err := r.db.First(&sample, sampleId).Error; err != nil {
		return err
	}

	var experiment ds.ImpurityFractionExperiment
	if err := r.db.First(&experiment, experimentId).Error; err != nil {
		return err
	}
	
	var existingLink ds.ExperimentSample
	err := r.db.Where("experiment_id = ? and sample_id = ?", experimentId, sampleId).First(&existingLink).Error
	if err == nil {
		// a link already exists
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// another error occurred
		return err
	}
	
	return r.db.Create(&ds.ExperimentSample{
		ExperimentID: experimentId,
		SampleID: sampleId,
	}).Error
}


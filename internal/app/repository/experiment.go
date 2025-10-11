package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetExperiments(from, to time.Time, status string) ([]ds.ImpurityFractionExperiment, error) {
	var experiments []ds.ImpurityFractionExperiment
	sub := r.db.Where("status != 'deleted' and status != 'draft'")
	if !from.IsZero() {
		sub = sub.Where("created_at > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("created_at < ?", to.Add(time.Hour*24))
	}

	if status != "" {
		sub = sub.Where("status = ?", status)
	}

	err := sub.Order("id").Find(&experiments).Error
	if err != nil {
		return nil, err
	}
	return experiments, nil
}

func (r *Repository) GetExperimentSamples(experimentId uint) ([]ds.ExperimentSample, error) {
	var experimentSamples []ds.ExperimentSample
	err := r.db.Where("experiment_id = ?", experimentId).Find(&experimentSamples).Error
	if err != nil {
		return nil, err
	}
	return experimentSamples, nil
}

func (r *Repository) GetExperimentSample(SampleID uint, ExperimentID uint) (ds.ExperimentSample, error) {
	var experimentSample ds.ExperimentSample
	err := r.db.Where("sample_id = ? and experiment_id = ?", SampleID, ExperimentID).First(&experimentSample).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ExperimentSample{}, fmt.Errorf("%w: образец эксперимента не найден", ErrNotFound)
		}
		return ds.ExperimentSample{}, err
	}
	return experimentSample, nil
}

func (r *Repository) GetExperimentSamplesData(id uint) ([]ds.AcidSolubleSample, ds.ImpurityFractionExperiment, error) {
	experiment, err := r.GetSingleExperiment(id)
	if err != nil {
		return []ds.AcidSolubleSample{}, ds.ImpurityFractionExperiment{}, err
	}

	var samples []ds.AcidSolubleSample
	sub := r.db.Table("experiments_samples").Where("experiment_id = ?", experiment.ID)
	err = r.db.Order("id DESC").Where("id IN (?)", sub.Select("sample_id")).Find(&samples).Error

	if err != nil {
		return []ds.AcidSolubleSample{}, ds.ImpurityFractionExperiment{}, err
	}

	return samples, experiment, nil
}

func (r *Repository) CheckCurrentExperimentDraft(creatorID uuid.UUID) (ds.ImpurityFractionExperiment, error) {	
	var experiment ds.ImpurityFractionExperiment
	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&experiment)
	if res.Error != nil {
		return ds.ImpurityFractionExperiment{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.ImpurityFractionExperiment{}, ErrNoDraft
	}
	return experiment, nil
}

func (r *Repository) GetExperimentDraft(creatorID uuid.UUID) (ds.ImpurityFractionExperiment, bool, error) {
	experiment, err := r.CheckCurrentExperimentDraft(creatorID)
	if errors.Is(err, ErrNoDraft) {
		experiment = ds.ImpurityFractionExperiment{
			Status:     "draft",
			CreatorID:  creatorID,
			CreatedAt: time.Now(),
		}
		result := r.db.Create(&experiment)
		if result.Error != nil {
			return ds.ImpurityFractionExperiment{}, false, result.Error
		}
		return experiment, true, nil
	} else if err != nil {
		return ds.ImpurityFractionExperiment{}, false, err
	}
	return experiment, true, nil
}

func (r *Repository) GetExperimentCount(creatorID uuid.UUID) int64 {
	if creatorID == uuid.Nil {
			return 0
	}
		
	var count int64
	experiment, err := r.CheckCurrentExperimentDraft(creatorID)
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.ExperimentSample{}).Where("experiment_id = ?", experiment.ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in experiments_samples:", err)
	}

	return count
}

func (r *Repository) GetSingleExperiment(id uint) (ds.ImpurityFractionExperiment, error) {
	var experiment ds.ImpurityFractionExperiment
	err := r.db.Where("id = ?", id).First(&experiment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ImpurityFractionExperiment{}, fmt.Errorf("%w: эксперимент с id %d", ErrNotFound, id)
		}
		return ds.ImpurityFractionExperiment{}, err
	} else if experiment.Status == "deleted"  {
		return ds.ImpurityFractionExperiment{}, fmt.Errorf("%w: эксперимент удален", ErrNotAllowed)
	}
	return experiment, nil
}

func (r *Repository) FormExperiment(experimentId uint, status string) (ds.ImpurityFractionExperiment, error) {
	experiment, err := r.GetSingleExperiment(experimentId)
	if err != nil {
		return ds.ImpurityFractionExperiment{}, err
	}

	if experiment.Status != "draft" {
		return ds.ImpurityFractionExperiment{}, fmt.Errorf("эта заявка не может быть %s", status)
	}
	
	if status != "deleted"{
		if experiment.MolarVolume == nil || *experiment.MolarVolume <= 0 {
		 	return ds.ImpurityFractionExperiment{}, errors.New("некорректный молярный объем")
		}
		experimentSamples, _ := r.GetExperimentSamples(experiment.ID)
		for _, experimentSample := range experimentSamples {
			if experimentSample.SampleMass == nil || *experimentSample.SampleMass <= 0{
				return ds.ImpurityFractionExperiment{}, errors.New("некорректная масса образца" )			
			}
			if experimentSample.EvolvedGasVolume == nil || *experimentSample.EvolvedGasVolume <= 0{
				return ds.ImpurityFractionExperiment{}, errors.New("некорректный объем выделившегося газа" )			
			}
		}
	}	

	now := time.Now()
	err = r.db.Model(&experiment).Updates(ds.ImpurityFractionExperiment{
		Status:   status,
		FormedAt: &now,
	}).Error
	if err != nil {
		return ds.ImpurityFractionExperiment{}, err
	}

	return experiment, nil
}

func (r *Repository) UpdateExperiment(id uint, experiment ds.ImpurityFractionExperiment) (ds.ImpurityFractionExperiment, error) {
	dbExperiment := ds.ImpurityFractionExperiment{}

	if experiment.MolarVolume != nil && *experiment.MolarVolume <= 0 {
		return ds.ImpurityFractionExperiment{}, errors.New("некорректный молярный объем")
  }

	err := r.db.Where("id = ? and status != 'deleted'", id).First(&dbExperiment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ImpurityFractionExperiment{}, fmt.Errorf("%w: эксперимент с id %d", ErrNotFound, id)
		}
		return ds.ImpurityFractionExperiment{}, err
	}
	err = r.db.Model(&dbExperiment).Updates(experiment).Error
	if err != nil {
		return ds.ImpurityFractionExperiment{}, err
	}
	return dbExperiment, nil
}

func CalculateMassFraction(
	molarVolume float64,
	relativeMolecularMass float64,
	stoichiometricCoefficient float64,
	sampleMass float64,
	evolvedGasVolume float64,
) (float64, error) {

	if molarVolume == 0 {
		return 0, errors.New("неправильный молярный объем")
	}

	if sampleMass == 0 {
		return 0, errors.New("неправильная масса образца")
	}
	
	if evolvedGasVolume == 0 {
		return 0, errors.New("неправильный объем выделившегося газа")
	}

	// Найдём количество вещества газа в молях
	nGas := evolvedGasVolume / molarVolume
	// Найдём количество вещества чистого вещества в молях
	nPureSubstance := nGas / stoichiometricCoefficient
	// Найдём массу чистого вещества в граммах
	mPureSubstance := nPureSubstance * relativeMolecularMass
	// Найдём массу примесей в граммах
	mImpurities := sampleMass - mPureSubstance
	// Найдём массовую долю примесей в процентах
	massFractionPercentage := mImpurities / sampleMass * 100
	return massFractionPercentage, nil
}

func (r *Repository) ModerateExperiment(id uint, status string, currUserId uuid.UUID) (ds.ImpurityFractionExperiment, error) {
	if status != "finished" && status != "rejected" {
		return ds.ImpurityFractionExperiment{}, errors.New("неверный статус")
	}

	experiment, err := r.GetSingleExperiment(id)
	if err != nil {
		return ds.ImpurityFractionExperiment{}, err
	} else if experiment.Status != "formed" {
		return ds.ImpurityFractionExperiment{}, errors.New("this experiment can not be " + status)
	}

	if status == "finished" {
		experimentSamples, err := r.GetExperimentSamples(experiment.ID)
		if err != nil {
			return ds.ImpurityFractionExperiment{}, err
		}
		for _, experimentSample := range experimentSamples {
			sample, err := r.GetSample(experimentSample.SampleID)
			if err != nil {
				return ds.ImpurityFractionExperiment{}, err
			}
			massFractionPercentage, err := CalculateMassFraction(*experiment.MolarVolume, sample.RelativeMolecularMass, sample.StoichiometricCoefficient, *experimentSample.SampleMass, *experimentSample.EvolvedGasVolume)
			if err != nil {
				return ds.ImpurityFractionExperiment{}, err
			}
			err = r.db.Model(&experimentSample).Where("experiment_id = ? AND sample_id = ?", experiment.ID, experimentSample.SampleID).Updates(ds.ExperimentSample{
				MassFractionPercentage: &massFractionPercentage,
			}).Error
			if err != nil {
				return ds.ImpurityFractionExperiment{}, err
			}
		}
	}

	now := time.Now()
	err = r.db.Model(&experiment).Updates(ds.ImpurityFractionExperiment{
		Status: status,
		FinishedAt: &now,
		ModeratorID: &currUserId,
	}).Error
	if err != nil {
		return ds.ImpurityFractionExperiment{}, err
	}

	return experiment, nil
}
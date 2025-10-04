package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/findwannawhy/IAD-Semester5/internal/app/minio"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (r *Repository) GetSamples() ([]ds.AcidSolubleSample, error) {
	var samples []ds.AcidSolubleSample
	err := r.db.Order("id").Where("deleted = false").Find(&samples).Error
	if err != nil {
		return nil, err
	}
	if len(samples) == 0 {
		return nil, fmt.Errorf("образцов не найдено")
	}

	return samples, nil
}

func (r *Repository) GetSample(id uint) (*ds.AcidSolubleSample, error) {
	sample := ds.AcidSolubleSample{}
	err := r.db.Order("id").Where("id = ? and deleted = ?", id, false).First(&sample).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: образец с id %d", ErrNotFound, id)
		}
		return &ds.AcidSolubleSample{}, err
	}
	return &sample, nil
}

func (r *Repository) GetSamplesByName(name string) ([]ds.AcidSolubleSample, error) {
	var samples []ds.AcidSolubleSample
	err := r.db.Order("id").Where("(title ILIKE ? OR formula ILIKE ?) AND deleted = ?", "%"+name+"%", "%"+name+"%", false).Find(&samples).Error
	if err != nil {
		return nil, err
	}
	return samples, nil
}

func (r *Repository) CreateSample(sample ds.AcidSolubleSample) (ds.AcidSolubleSample, error) {
	if sample.RelativeMolecularMass <= 0 {
		return ds.AcidSolubleSample{}, errors.New("некорректная относительная молекулярная масса")
	}
	if sample.StoichiometricCoefficient <= 0 {
		return ds.AcidSolubleSample{}, errors.New("некорректный стехиометрический коэффициент")
	}
	err := r.db.Create(&sample).Error
	if err != nil {
		return ds.AcidSolubleSample{}, err
	}
	return sample, nil
}

func (r *Repository) UpdateSample(id uint, sample ds.AcidSolubleSample) (ds.AcidSolubleSample, error) {
	dbSample := ds.AcidSolubleSample{}
	err := r.db.Where("id = ? and deleted = ?", id, false).First(&dbSample).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.AcidSolubleSample{}, fmt.Errorf("%w: образец с id %d", ErrNotFound, id)
		}
		return ds.AcidSolubleSample{}, err
	}
	if sample.RelativeMolecularMass <= 0 {
		return ds.AcidSolubleSample{}, errors.New("некорректная относительная молекулярная масса")
	}
	if sample.StoichiometricCoefficient <= 0 {
		return ds.AcidSolubleSample{}, errors.New("некорректный стехиометрический коэффициент")
	}
	err = r.db.Model(&dbSample).Updates(sample).Error
	if err != nil {
		return ds.AcidSolubleSample{}, err
	}
	return dbSample, nil
}

func (r *Repository) SoftDeleteSample(id uint) error {
	sample := ds.AcidSolubleSample{}

	err := r.db.Where("id = ? and deleted = ?", id, false).First(&sample).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: образец с id %d", ErrNotFound, id)
		}
		return err
	}
	if sample.ImageURL != nil {
		err = minio.DeleteObject(context.Background(), r.mc, minio.GetImgBucket(), *sample.ImageURL)
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.AcidSolubleSample{}).Where("id = ?", id).Update("deleted", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddSampleToExperimentDraft(experimentId uint, sampleId uint) error {
	var sample ds.AcidSolubleSample
	if err := r.db.First(&sample, sampleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: образец с id %d", ErrNotFound, sampleId)
		}
		return err
	}

	var experiment ds.ImpurityFractionExperiment
	if err := r.db.First(&experiment, experimentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: эксперимент с id %d", ErrNotFound, experimentId)
		}
		return err
	}
	
	experimentSample := ds.ExperimentSample{}
	result := r.db.Where("sample_id = ? and experiment_id = ?", sampleId, experimentId).Find(&experimentSample)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return fmt.Errorf("%w: образец %d уже в эксперименте %d", ErrAlreadyExists, sampleId, experimentId)
	}
	return r.db.Create(&ds.ExperimentSample{
		SampleID:     uint(sampleId),
		ExperimentID: uint(experimentId),
	}).Error
}

func (r *Repository) GetModeratorAndCreatorLogin(experiment ds.ImpurityFractionExperiment) (string, string, error) {
	var creator ds.User
	var moderator ds.User

	err := r.db.Where("id = ?", experiment.CreatorID).First(&creator).Error
	if err != nil {
		return "", "", err
	}

	var moderatorLogin string
	if experiment.ModeratorID != nil {
		err = r.db.Where("id = ?", *experiment.ModeratorID).First(&moderator).Error
		if err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}
	
	return creator.Login, moderatorLogin, nil
}

func (r *Repository) UpdateImage(ctx *gin.Context, sampleId uint, file *multipart.FileHeader) (ds.AcidSolubleSample, error) {
	sample_, err := r.GetSample(sampleId)
	if err != nil {
		return ds.AcidSolubleSample{}, err
	}
	
	fileName, err := minio.UploadImage(ctx, r.mc, minio.GetImgBucket(), file, *sample_)
	if err != nil {
		return ds.AcidSolubleSample{},err
	}

	sample, err := r.GetSample(sampleId)
	if err != nil {
		return ds.AcidSolubleSample{}, err
	}
	sample.ImageURL = &fileName
	err = r.db.Save(&sample).Error
	if err != nil {
		return ds.AcidSolubleSample{}, err
	}
	return *sample, nil
}


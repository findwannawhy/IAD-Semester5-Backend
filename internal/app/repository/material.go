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

func (r *Repository) GetMaterials() ([]ds.Material, error) {
	var materials []ds.Material
	err := r.db.Order("id").Where("deleted = false").Find(&materials).Error
	if err != nil {
		return nil, err
	}
	if len(materials) == 0 {
		return nil, fmt.Errorf("веществ не найдено")
	}

	return materials, nil
}

func (r *Repository) GetMaterial(id uint) (*ds.Material, error) {
	material := ds.Material{}
	err := r.db.Order("id").Where("id = ? and deleted = ?", id, false).First(&material).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: материал с id %d", ErrNotFound, id)
		}
		return &ds.Material{}, err
	}
	return &material, nil
}

func (r *Repository) GetMaterialsByName(name string) ([]ds.Material, error) {
	var materials []ds.Material
	err := r.db.Order("id").Where("(title ILIKE ? OR formula ILIKE ?) AND deleted = ?", "%"+name+"%", "%"+name+"%", false).Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}

func (r *Repository) CreateMaterial(material ds.Material) (ds.Material, error) {
	if material.RelativeMolecularMass <= 0 {
		return ds.Material{}, errors.New("некорректная относительная молекулярная масса")
	}
	if material.StoichiometricCoefficient <= 0 {
		return ds.Material{}, errors.New("некорректный стехиометрический коэффициент")
	}
	err := r.db.Create(&material).Error
	if err != nil {
		return ds.Material{}, err
	}
	return material, nil
}

func (r *Repository) UpdateMaterial(id uint, material ds.Material) (ds.Material, error) {
	dbMaterial := ds.Material{}
	err := r.db.Where("id = ? and deleted = ?", id, false).First(&dbMaterial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Material{}, fmt.Errorf("%w: вещество с id %d", ErrNotFound, id)
		}
		return ds.Material{}, err
	}
	if material.RelativeMolecularMass <= 0 {
		return ds.Material{}, errors.New("некорректная относительная молекулярная масса")
	}
	if material.StoichiometricCoefficient <= 0 {
		return ds.Material{}, errors.New("некорректный стехиометрический коэффициент")
	}
	err = r.db.Model(&dbMaterial).Updates(material).Error
	if err != nil {
		return ds.Material{}, err
	}
	return dbMaterial, nil
}

func (r *Repository) SoftDeleteMaterial(id uint) error {
	material := ds.Material{}

	err := r.db.Where("id = ? and deleted = ?", id, false).First(&material).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: вещество с id %d", ErrNotFound, id)
		}
		return err
	}
	if material.ImageURL != nil {
		err = minio.DeleteObject(context.Background(), r.mc, minio.GetImgBucket(), *material.ImageURL)
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.Material{}).Where("id = ?", id).Update("deleted", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddMaterialToExperiment(experimentId uint, materialId uint) error {
	var material ds.Material
	if err := r.db.First(&material, materialId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: вещество с id %d", ErrNotFound, materialId)
		}
		return err
	}

	var experiment ds.Experiment
	if err := r.db.First(&experiment, experimentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: исследование с id %d", ErrNotFound, experimentId)
		}
		return err
	}
	
	experimentItem := ds.ExperimentItem{}
	result := r.db.Where("material_id = ? and experiment_id = ?", materialId, experimentId).Find(&experimentItem)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return fmt.Errorf("%w: вещество %d уже в исследованиии %d", ErrAlreadyExists, materialId, experimentId)
	}
	return r.db.Create(&ds.ExperimentItem{
		MaterialID:    uint(materialId),
		ExperimentID: uint(experimentId),
	}).Error
}

func (r *Repository) GetModeratorAndCreatorLogin(experiment ds.Experiment) (string, string, error) {
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

func (r *Repository) UpdateImage(ctx *gin.Context, materialId uint, file *multipart.FileHeader) (ds.Material, error) {
	material_, err := r.GetMaterial(materialId)
	if err != nil {
		return ds.Material{}, err
	}
	
	fileName, err := minio.UploadImage(ctx, r.mc, minio.GetImgBucket(), file, *material_)
	if err != nil {
		return ds.Material{},err
	}

	material, err := r.GetMaterial(materialId)
	if err != nil {
		return ds.Material{}, err
	}
	material.ImageURL = &fileName
	err = r.db.Save(&material).Error
	if err != nil {
		return ds.Material{}, err
	}
	return *material, nil
}
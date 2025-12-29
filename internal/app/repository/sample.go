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

// PaginatedSamples содержит образцы с метаданными пагинации
type PaginatedSamples struct {
	Samples    []ds.AcidSolubleSample `json:"samples"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

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

// GetSamplesPaginated возвращает образцы с пагинацией
func (r *Repository) GetSamplesPaginated(name string, page, limit int) (*PaginatedSamples, error) {
	var samples []ds.AcidSolubleSample
	var total int64

	// Базовый запрос
	query := r.db.Model(&ds.AcidSolubleSample{}).Where("deleted = false")

	// Фильтрация по имени (если задано)
	if name != "" {
		query = query.Where("title ILIKE ? OR formula ILIKE ?", "%"+name+"%", "%"+name+"%")
	}

	// Подсчет общего количества
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Вычисление offset
	offset := (page - 1) * limit

	// Получение данных с пагинацией
	if err := query.Order("id").Limit(limit).Offset(offset).Find(&samples).Error; err != nil {
		return nil, err
	}

	// Вычисление общего количества страниц
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &PaginatedSamples{
		Samples:    samples,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
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

// GetSamplesByIDs возвращает образцы по списку ID в нужном порядке
func (r *Repository) GetSamplesByIDs(ids []uint) ([]ds.AcidSolubleSample, error) {
	if len(ids) == 0 {
		return []ds.AcidSolubleSample{}, nil
	}
	
	var samples []ds.AcidSolubleSample
	err := r.db.Where("id IN ? AND deleted = ?", ids, false).Find(&samples).Error
	if err != nil {
		return nil, err
	}
	
	// Сортируем образцы в соответствии с порядком ID
	sampleMap := make(map[uint]ds.AcidSolubleSample)
	for _, sample := range samples {
		sampleMap[sample.ID] = sample
	}
	
	orderedSamples := make([]ds.AcidSolubleSample, 0, len(ids))
	for _, id := range ids {
		if sample, ok := sampleMap[id]; ok {
			orderedSamples = append(orderedSamples, sample)
		}
	}
	
	return orderedSamples, nil
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


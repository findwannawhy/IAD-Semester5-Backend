package ds

import (
	"time"
)

type Material struct {
	ID                         int        `gorm:"primaryKey"                json:"id"`
	Title                      string     `gorm:"size:50;not null;index"    json:"title"`
	Formula                    string     `gorm:"size:30;not null;index"    json:"formula"`
	Description               *string     `gorm:"type:text"                 json:"description"`
	Deleted                    bool       `gorm:"not null;default:false"    json:"deleted"`
	ImageURL                  *string     `gorm:"size:255"                  json:"image_url"`
	RelativeMolecularMass      float64    `gorm:"not null"                  json:"relative_molecular_mass"`
	StoichiometricCoefficient  float64    `gorm:"not null"                  json:"stoichiometric_coefficient"`
	
	CreatedAt                  time.Time  `gorm:"<-:create;not null"        json:"created_at"`
	UpdatedAt                  time.Time  `gorm:"not null"                  json:"updated_at"`
}
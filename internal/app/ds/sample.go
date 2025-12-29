package ds

type AcidSolubleSample struct {
	ID                        uint    `gorm:"primaryKey;index:idx_samples_pagination,priority:2"  json:"id"`
	Title                     string  `gorm:"size:50;not null;index"                              json:"title"`
	Formula                   string  `gorm:"size:30;not null;index"                              json:"formula"`
	Description               *string `gorm:"type:text"                                           json:"description"`
	Deleted                   bool    `gorm:"not null;default:false;index:idx_samples_pagination,priority:1" json:"deleted"`
	ImageURL                  *string `gorm:"size:255"                                            json:"image_url"`
	RelativeMolecularMass     float64 `gorm:"type:double precision;not null"                      json:"relative_molecular_mass"`
	StoichiometricCoefficient float64 `gorm:"type:double precision;not null"                      json:"stoichiometric_coefficient"`
}
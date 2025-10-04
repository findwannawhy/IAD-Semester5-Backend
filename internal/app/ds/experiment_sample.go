package ds

type ExperimentSample struct {
	ExperimentID             uint       `gorm:"primaryKey;autoIncrement:false" json:"experiment_id"`
	SampleID                 uint       `gorm:"primaryKey;autoIncrement:false" json:"sample_id"`

	SampleMass              *float64    `gorm:"type:double precision"          json:"sample_mass"`
	EvolvedGasVolume        *float64    `gorm:"type:double precision"          json:"evolved_gas_volume"`
	MassFractionPercentage  *float64    `gorm:"type:double precision"          json:"mass_fraction_percentage"`

	Experiment              *ImpurityFractionExperiment `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"   json:"impurity_fraction_experiment"`
	Sample                  *AcidSolubleSample          `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"   json:"acid_soluble_sample"`
}

func (ExperimentSample) TableName() string {
	return "experiments_samples"
}
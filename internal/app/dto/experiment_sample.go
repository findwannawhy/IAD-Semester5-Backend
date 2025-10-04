package dto

type DeleteExperimentSampleResp struct {
	ExperimentID   uint       `json:"experiment_id"`
	SampleID       uint       `json:"sample_id"`
	CreatorLogin   string     `json:"creator_login"`
}

type UpdateExperimentSampleReq struct {
	SampleMass       *float64 `json:"sample_mass"`
	EvolvedGasVolume *float64 `json:"evolved_gas_volume"`
}
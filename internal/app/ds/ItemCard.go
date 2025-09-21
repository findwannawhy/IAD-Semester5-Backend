package ds

type ItemCard struct {
	ItemID                     int
	Title                      string
	Description               *string
	ImageURL                  *string
	RelativeMolecularMass      float64
	StoichiometricCoefficient  float64
	MaterialMass               string
	GasVolume                  string
	MassFractionPercentage     string
}
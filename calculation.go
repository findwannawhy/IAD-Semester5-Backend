package main

import "fmt"

type CalculateMassFractionInput struct {
	MolarVolume               float64 // может быть не задан (0.0) (поле заявки - заполняет пользователь)
	RelativeMolecularMass     float64 // обязан быть задан, нельзя не задать (админ создает карточку)
	StoichiometricCoefficient float64 // обязан быть задан, нельзя не задать (админ создает карточку)
	SampleMass              float64 // обязан быть задан, если пользователь оставил 0.0, выбрасываем 0 (м-м поле - заполняет пользователь)
	GasVolume                 float64 // обязан быть задан, если пользователь оставил 0.0, выбрасываем 0 (м-м поле - заполняет пользователь)
}

func (input CalculateMassFractionInput) CalculateMassFraction() float64 {

	if input.MolarVolume == 0 {
		input.MolarVolume = 22.4
	}

	if input.SampleMass == 0 || input.GasVolume == 0 {
		return 0
	}

	// Найдём количество вещества газа в молях
	nGas := input.GasVolume / input.MolarVolume
	// Найдём количество вещества чистого вещества в молях
	nPureSubstance := nGas / input.StoichiometricCoefficient
	// Найдём массу чистого вещества в граммах
	mPureSubstance := nPureSubstance * input.RelativeMolecularMass
	// Найдём массу примесей в граммах
	mImpurities := input.SampleMass - mPureSubstance
	// Найдём массовую долю примесей в процентах
	massFractionPercentage := mImpurities / input.SampleMass * 100
	return massFractionPercentage
}

func main() {
	input := CalculateMassFractionInput{
		MolarVolume:               22.4,   // Молярный объем газа (л/моль)
		RelativeMolecularMass:     100.07, // Относительная молекулярная масса вещества (г/моль)
		StoichiometricCoefficient: 1.0,    // Стехиометрический коэффициент
		SampleMass:              2.35,   // Масса исходного образца материала (г)
		GasVolume:                 0.484,  // Объем выделившегося газа (л)
	}

	massFractionPercentage := input.CalculateMassFraction()

	fmt.Printf("Массовая доля примесей: %.2f%%\n", massFractionPercentage)
}

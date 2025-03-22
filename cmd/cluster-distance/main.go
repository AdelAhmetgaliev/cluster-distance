package main

import (
	"fmt"
	"log"
	"math"
	"path/filepath"

	"github.com/AdelAhmetgaliev/cluster-distance/internal/stardata"
	"github.com/AdelAhmetgaliev/cluster-distance/internal/utils"
	"gonum.org/v1/gonum/interp"
)

func main() {
	// Считаем входные данные звезд скопления
	inputStarsFilePath := filepath.Join("data", "input", "stars_NGC-869.csv")
	inputStarsSlice := utils.ReadStars(inputStarsFilePath)

	// Из всех звезд выберем малую группу, которую можно совместить с линией нормальных цветов
	// Наилучшим образом подходят звезды лежащие по следующим координатам: BV = [0.2; 0.4] UB = [-0.6; -0.2]
	var processingStarsSlice []stardata.StarData
	for _, sd := range inputStarsSlice {
		if sd.CI.BV >= 0.2 && sd.CI.BV <= 0.4 && sd.CI.UB >= -0.6 && sd.CI.UB <= -0.2 {
			processingStarsSlice = append(processingStarsSlice, sd)
		}
	}

	// Считаем входные показатели цвета главной последовательности
	inputColorIndexesFilePath := filepath.Join("data", "input", "color_indexes.csv")
	inputColorIndexesSlice := utils.ReadColorIndexes(inputColorIndexesFilePath)

	// Интерполируем входные показатели цвета. Интерполяцию будем делать Акимовскими сплайнами
	var akimaInterpOfInputColorIndexes interp.AkimaSpline
	{
		// Разделим данные на два массива
		sliceCap := len(inputColorIndexesSlice)
		xValues := make([]float64, 0, sliceCap)
		yValues := make([]float64, 0, sliceCap)
		for _, ci := range inputColorIndexesSlice {
			xValues = append(xValues, ci.BV)
			yValues = append(yValues, ci.UB)
		}

		if err := akimaInterpOfInputColorIndexes.Fit(xValues, yValues); err != nil {
			log.Fatalf("Can't interpolate input color indexes: %v\n", err)
		}
	}

	// Найдем средний показатель цвета обрабатываемых звезд
	processingStarsAverageColorIndex := stardata.AverageColorIndexes(processingStarsSlice)

	// Найдем уравнение линии покраснения для среднего показателя цвета
	// Уравнение имеет вид y = k0 + K * x, где y - UB, x - BV, K = 0.72
	const K = 0.72
	k0 := processingStarsAverageColorIndex.UB - K*processingStarsAverageColorIndex.BV

	// Найдем минимальное и максимальное значение BV входных показателей цвета главной последовательности
	bvMin, bvMax := math.Inf(1), math.Inf(-1)
	for _, ci := range inputColorIndexesSlice {
		if ci.BV < bvMin {
			bvMin = ci.BV
		}
		if ci.BV > bvMax {
			bvMax = ci.BV
		}
	}

	// Найдем пересечение линии покраснения с линией нормальных цветов
	bvIntersec, ubIntersec := math.Inf(1), math.Inf(1)
	for bv := bvMin - 0.3; bv <= processingStarsAverageColorIndex.BV; bv += 0.0001 {
		ub := k0 + K*bv
		if math.Abs(ub-akimaInterpOfInputColorIndexes.Predict(bv)) <= 0.001 {
			bvIntersec = bv
			ubIntersec = ub
			break
		}
	}

	if math.IsInf(bvIntersec, 1) || math.IsInf(ubIntersec, 1) {
		log.Fatalln("Failed to find the intersection point")
	}

	// Найдем смещение среднего показателя цвета
	bvDelta := bvIntersec - processingStarsAverageColorIndex.BV
	ubDelta := ubIntersec - processingStarsAverageColorIndex.UB

	// Сместим все обрабатываемые звезды
	correctedStarsSlice := make([]stardata.StarData, 0, len(processingStarsSlice))
	for _, sd := range processingStarsSlice {
		correctedStarData := sd
		correctedStarData.CI.BV += bvDelta
		correctedStarData.CI.UB += ubDelta

		correctedStarsSlice = append(correctedStarsSlice, correctedStarData)
	}

	// Считаем значения звездной величины звезд ГП
	inputMagVToBVFilePath := filepath.Join("data", "input", "color_indexes.csv")
	inputMagVToBV := utils.ReadMagVToBV(inputMagVToBVFilePath)

	// Интерполируем значения звездной величины звезд ГП
	var akimaInterpOfInputMagVToBV interp.AkimaSpline
	{
		// Разделим данные на два массива
		sliceCap := len(inputMagVToBV)
		xValues := make([]float64, 0, sliceCap)
		yValues := make([]float64, 0, sliceCap)
		for _, chunk := range inputMagVToBV {
			xValues = append(xValues, chunk[0])
			yValues = append(yValues, chunk[1])
		}

		if err := akimaInterpOfInputMagVToBV.Fit(xValues, yValues); err != nil {
			log.Fatalf("Can't interpolate input mag v: %v\n", err)
		}
	}

	// Рассчитаем среднее значение звездной величины в фильтре V обрабатываемых звезд
	correctedStarsAverageMagV := 0.0
	for _, sd := range processingStarsSlice {
		correctedStarsAverageMagV += sd.Mag.V
	}
	correctedStarsAverageMagV /= float64(len(processingStarsSlice))

	// Найдем значение пересечения средней звезды с линией звезд ГП
	magVIntersec := akimaInterpOfInputMagVToBV.Predict(bvIntersec)

	// Сместим обрабатываемые звезды
	magVDelta := magVIntersec - correctedStarsAverageMagV
	for i := range correctedStarsSlice {
		correctedStarsSlice[i].Mag.V += magVDelta
	}

	// Найдем расстояние до скопления: mv - Mv = 5 * lg(r) - 5 + Rv * E(B-V)
	const Rv = 3.1
	colorExcess := -bvDelta
	deltaMagV := -magVDelta
	distance := math.Pow(10, (deltaMagV+5.0-Rv*colorExcess)/5.0)

	fmt.Printf("Расстояние до скопления: %.1f пк\n", distance)

	// Выведем полученные данные в отдельные файлы
	outputDirFilePath := filepath.Join("data", "output")

	processingStarsColorIndexesFilePath := filepath.Join(outputDirFilePath, "processing_stars_color_indexes.dat")
	stardata.WriteColorIndexesToFile(processingStarsColorIndexesFilePath, processingStarsSlice)

	processingStarsMagVToBVFilePath := filepath.Join(outputDirFilePath, "processing_stars_magv_to_bv.dat")
	stardata.WriteMagVToBVToFile(processingStarsMagVToBVFilePath, processingStarsSlice)

	correctedStarsColorIndexesFilePath := filepath.Join(outputDirFilePath, "corrected_stars_color_indexes.dat")
	stardata.WriteColorIndexesToFile(correctedStarsColorIndexesFilePath, correctedStarsSlice)

	correctedStarsMagVToBVFilePath := filepath.Join(outputDirFilePath, "corrected_stars_magv_to_bv.dat")
	stardata.WriteMagVToBVToFile(correctedStarsMagVToBVFilePath, correctedStarsSlice)

	averageColorIndexFilePath := filepath.Join(outputDirFilePath, "average_color_index.dat")
	processingStarsAverageColorIndex.WriteToFile(averageColorIndexFilePath)

	correctedColorIndexFilePath := filepath.Join(outputDirFilePath, "corrected_color_index.dat")
	correctedColorIndex := stardata.NewColorIndex(bvIntersec, ubIntersec)
	correctedColorIndex.WriteToFile(correctedColorIndexFilePath)

	// Заполним срезы интерполированными данными для рисования графиков
	bvStep := 0.01
	inputColorIndexesInterpolated := make([][2]float64, 0, int((bvMax-bvMin)/bvStep)+5)
	inputMagVToBVInterpolated := make([][2]float64, 0, int((bvMax-bvMin)/bvStep)+5)
	for bv := bvMin; bv <= bvMax; bv += bvStep {
		var chunkCI [2]float64
		chunkCI[0] = bv
		chunkCI[1] = akimaInterpOfInputColorIndexes.Predict(bv)

		var chunkMagVToBV [2]float64
		chunkMagVToBV[0] = bv
		chunkMagVToBV[1] = akimaInterpOfInputMagVToBV.Predict(bv)

		inputColorIndexesInterpolated = append(inputColorIndexesInterpolated, chunkCI)
		inputMagVToBVInterpolated = append(inputMagVToBVInterpolated, chunkMagVToBV)
	}

	inputColorIndexesInterpolatedFilePath := filepath.Join(outputDirFilePath, "color_indexes_interpolated.dat")
	utils.WriteSliceToFile(inputColorIndexesInterpolatedFilePath, inputColorIndexesInterpolated)

	inputMagVToBVInterpolatedFilePath := filepath.Join(outputDirFilePath, "magv_to_bv_interpolated.dat")
	utils.WriteSliceToFile(inputMagVToBVInterpolatedFilePath, inputMagVToBVInterpolated)
}

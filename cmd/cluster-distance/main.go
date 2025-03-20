package main

import (
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/AdelAhmetgaliev/cluster-distance/internal/stardata"
	"github.com/AdelAhmetgaliev/cluster-distance/internal/utils"
	"gonum.org/v1/gonum/interp"
)

func main() {
	// Считаем исходные данные из csv-файлов
	starsList := utils.ReadStars("data/stars_NGC-869.csv")
	mainColorIndexes := utils.ReadColorIndexes("data/bolometric_corrections_V.csv")

	// Определим минимальные и максимальные значения показателей цвета
	bvMax, bvMin := -9999.0, 9999.0
	ubMax, ubMin := -9999.0, 9999.0
	for _, ci := range mainColorIndexes {
		if bvMax < ci[0] {
			bvMax = ci[0]
		}
		if bvMin > ci[0] {
			bvMin = ci[0]
		}
		if ubMax < ci[1] {
			ubMax = ci[1]
		}
		if ubMin > ci[1] {
			ubMin = ci[1]
		}
	}

	// Исключим из обработки звезды, которые сильно выбиваются своими показателями цвета
	processedStarsList := make([]stardata.StarData, 0, len(starsList))
	for _, sd := range starsList {
		var ci [2]float64
		ci[0], ci[1] = sd.CI.BV, sd.CI.UB

		const delta = 0.3
		bv, ub := ci[0], ci[1]
		if bv > bvMax && bv-bvMax > delta {
			continue
		}
		if bv < bvMin && bvMin-bv > delta {
			continue
		}

		if ub > ubMax && ub-ubMax > delta {
			continue
		}
		if ub < ubMin && ubMin-ub > delta {
			continue
		}

		processedStarsList = append(processedStarsList, sd)
	}
	utils.WriteSliceToFile("data/main_color_indexes.dat", mainColorIndexes)
	stardata.WriteColorIndexesToFile("data/stars_color_indexes.dat", processedStarsList)

	// Выделим три области звезд для их усреднения
	// Первая область BV > 0.8 UB > 0.0
	// Вторая область BV < 0.8 UB > 0.0
	// Третья область BV < 0.8 UB < 0.0
	var firstSetOfStars, secondSetOfStars, thirdSetOfStars []stardata.StarData
	for _, sd := range processedStarsList {
		if sd.CI.BV > 0.8 {
			firstSetOfStars = append(firstSetOfStars, sd)
			continue
		}

		if sd.CI.UB > 0.0 {
			secondSetOfStars = append(secondSetOfStars, sd)
			continue
		}

		thirdSetOfStars = append(thirdSetOfStars, sd)
	}
	stardata.WriteColorIndexesToFile("data/stars1_color_indexes.dat", firstSetOfStars)
	stardata.WriteColorIndexesToFile("data/stars2_color_indexes.dat", secondSetOfStars)
	stardata.WriteColorIndexesToFile("data/stars3_color_indexes.dat", thirdSetOfStars)

	// Усредним каждое из множеств
	averageCIOfFirstSet := stardata.AverageColorIndexes(firstSetOfStars)
	averageCIOfSecondSet := stardata.AverageColorIndexes(secondSetOfStars)
	averageCIOfThirdSet := stardata.AverageColorIndexes(thirdSetOfStars)
	averageCIOfFirstSet.WriteToFile("data/stars1_average_color_index.dat")
	averageCIOfSecondSet.WriteToFile("data/stars2_average_color_index.dat")
	averageCIOfThirdSet.WriteToFile("data/stars3_average_color_index.dat")

	// Интерполируем линию нормальных цветов
	var akimaInterp interp.AkimaSpline
	{
		var xValues, yValues []float64
		for _, ci := range mainColorIndexes {
			xValues = append(xValues, ci[0])
			yValues = append(yValues, ci[1])
		}

		if err := akimaInterp.Fit(xValues, yValues); err != nil {
			log.Fatalf("Can't interpolate: %v\n", err)
		}
	}

	var mainColorIndexesInterp [][2]float64
	for x := bvMin; x <= bvMax; x += 0.01 {
		var temp [2]float64
		temp[0] = x
		temp[1] = akimaInterp.Predict(x)

		mainColorIndexesInterp = append(mainColorIndexesInterp, temp)
	}
	utils.WriteSliceToFile("data/main_color_indexes_interp.dat", mainColorIndexesInterp)

	// Найдем пересечение с линией нормальных цветов для каждой звезды
	canBeCorrectedStarsList := make([]stardata.StarData, 0, len(processedStarsList))
	var correctedColorIndexes [][2]float64
	for _, sd := range processedStarsList {
		// Найдем уравнение линии покраснения: y[U - B] = k0 + K * x[B - V]
		const K = 0.72 // Наклон линии покраснения
		k0 := sd.CI.UB - K*sd.CI.BV

		bvIntersec := -100.0
		for bv := bvMin - 1; bv <= sd.CI.BV; bv += 0.0001 {
			if math.Abs(akimaInterp.Predict(bv)-(k0+K*bv)) < 0.01 {
				bvIntersec = bv
			}
		}
		// Если не нашли пересечение с линией нормальных цветов
		if bvIntersec == -100.0 {
			continue
		}

		var correctedCI [2]float64
		correctedCI[0] = bvIntersec
		correctedCI[1] = k0 + K*bvIntersec
		correctedColorIndexes = append(correctedColorIndexes, correctedCI)

		canBeCorrectedStarsList = append(canBeCorrectedStarsList, sd)
	}

	utils.WriteSliceToFile("data/stars_color_indexes_corrected.dat", correctedColorIndexes)
	stardata.WriteColorIndexesToFile("data/stars_color_indexes_can_be_corrected.dat", canBeCorrectedStarsList)

	magVToBV := utils.ReadDefaultMagVToBV("data/bolometric_corrections_V.csv")
	utils.WriteSliceToFile("data/main_magv_to_bv.dat", magVToBV)
	stardata.WriteMagVToBVToFile("data/stars_magv_to_bv.dat", canBeCorrectedStarsList)

	// Интерполируем ГР диаграмму для ГП
	{
		var xValues, yValues []float64
		for _, ci := range magVToBV {
			xValues = append(xValues, ci[0])
			yValues = append(yValues, ci[1])
		}

		if err := akimaInterp.Fit(xValues, yValues); err != nil {
			log.Fatalf("Can't interpolate: %v\n", err)
		}
	}

	var mainMagVToBvInterp [][2]float64
	for x := bvMin; x <= bvMax; x += 0.01 {
		var temp [2]float64
		temp[0] = x
		temp[1] = akimaInterp.Predict(x)

		mainMagVToBvInterp = append(mainMagVToBvInterp, temp)
	}
	utils.WriteSliceToFile("data/main_magv_to_bv_interp.dat", mainMagVToBvInterp)

	// Сделаем список откорректированных звезд по показателю цвета
	correctedStarsList := make([]stardata.StarData, 0, len(canBeCorrectedStarsList))
	for i, sd := range canBeCorrectedStarsList {
		correctedSD := sd
		correctedCI := stardata.NewColorIndex(correctedColorIndexes[i][0], correctedColorIndexes[i][1])
		correctedSD.CI = correctedCI

		correctedStarsList = append(correctedStarsList, correctedSD)
	}
	stardata.WriteMagVToBVToFile("data/stars_magv_to_bv_corrected.dat", correctedStarsList)

	// Рассчитаем скорректированные значения звездной величины в фильтре V
	correctedMagVList := make([]float64, 0, len(correctedStarsList))
	for _, sd := range correctedStarsList {
		correctedMagV := akimaInterp.Predict(sd.CI.BV)
		correctedMagVList = append(correctedMagVList, correctedMagV)
	}

	// Рассчитаем расстояния до звезд
	starDistanceList := make([]float64, 0, len(correctedStarsList))
	for i, sd := range canBeCorrectedStarsList {
		// mv - Mv = 5 * lg(r) - 5 + Rv * E(B-V)
		const RV = 3.1 // Rv

		excessColor := sd.CI.BV - correctedColorIndexes[i][0]
		deltaMag := sd.Mag.V - correctedMagVList[i]
		starDistance := math.Pow(10.0, (deltaMag+5.0-RV*excessColor)/5.0)

		starDistanceList = append(starDistanceList, starDistance)
	}

	sort.Slice(starDistanceList, func(i, j int) bool {
		return starDistanceList[i] < starDistanceList[j]
	})

	averageDistance := 0.0
	for _, d := range starDistanceList {
		averageDistance += d
	}
	averageDistance /= float64(len(starDistanceList))

	fmt.Printf("Среднее расстояние до РС:\t%.1f\n", averageDistance)
	fmt.Printf("Минимальное расстояние до РС:\t%.1f\n", starDistanceList[0])
	fmt.Printf("Максимальное расстояние до РС:\t%.1f\n", starDistanceList[len(starDistanceList)-1])
}

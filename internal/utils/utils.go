package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/AdelAhmetgaliev/cluster-distance/internal/stardata"
)

func ReadStars(starsFilePath string) []stardata.StarData {
	starsFile, err := os.Open(starsFilePath)
	if err != nil {
		log.Fatalf("Can't open file: %v\n", err)
	}

	defer func() {
		if err := starsFile.Close(); err != nil {
			log.Printf("Can't close file: %v\n", err)
		}
	}()

	reader := csv.NewReader(starsFile)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Can't read data from csv: %v\n", err)
	}

	starsList := stardata.CreateStarsList(data)

	return starsList
}

func ReadColorIndexes(colorIndexesFilePath string) [][2]float64 {
	file, err := os.Open(colorIndexesFilePath)
	if err != nil {
		log.Fatalf("Can't open file: %v\n", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Can't close file: %v\n", err)
		}
	}()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Can't read data from csv: %v\n", err)
	}

	colorIndexes := createDefaultColorIndexes(data)
	return colorIndexes
}

func createDefaultColorIndexes(data [][]string) (colorIndexes [][2]float64) {
	bvColumn, ubColumn := -1, -1
	for j, field := range data[0] {
		clearField := strings.TrimSpace(field)
		switch clearField {
		case "(B - V)0":
			bvColumn = j
		case "(U - B)0":
			ubColumn = j
		}
	}

	for i, line := range data {
		if i == 0 {
			continue
		}

		var colorIndex [2]float64
		for j, field := range line {
			clearField := strings.TrimSpace(field)
			switch j {
			case bvColumn:
				value, err := strconv.ParseFloat(clearField, 64)
				if err != nil {
					break
				}
				colorIndex[0] = value
			case ubColumn:
				value, err := strconv.ParseFloat(clearField, 64)
				if err != nil {
					break
				}
				colorIndex[1] = value
			}
		}

		colorIndexes = append(colorIndexes, colorIndex)
	}

	return colorIndexes
}

func ReadDefaultMagVToBV(filePath string) [][2]float64 {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Can't open file: %v\n", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Can't close file: %v\n", err)
		}
	}()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Can't read data from csv: %v\n", err)
	}

	magVToBV := createDefaultMagVToBV(data)
	return magVToBV
}

func createDefaultMagVToBV(data [][]string) (magVToBV [][2]float64) {
	bvColumn, magVColumn := -1, -1
	for j, field := range data[0] {
		clearField := strings.TrimSpace(field)
		switch clearField {
		case "(B - V)0":
			bvColumn = j
		case "MV":
			magVColumn = j
		}
	}

	for i, line := range data {
		if i == 0 {
			continue
		}

		var vToBV [2]float64
		for j, field := range line {
			clearField := strings.TrimSpace(field)
			switch j {
			case bvColumn:
				value, err := strconv.ParseFloat(clearField, 64)
				if err != nil {
					break
				}
				vToBV[0] = value
			case magVColumn:
				value, err := strconv.ParseFloat(clearField, 64)
				if err != nil {
					break
				}
				vToBV[1] = value
			}
		}

		magVToBV = append(magVToBV, vToBV)
	}

	return magVToBV
}

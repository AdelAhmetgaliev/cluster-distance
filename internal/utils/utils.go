package utils

import (
	"encoding/csv"
	"fmt"
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

	starsList := createStarsSlice(data)

	return starsList
}

func ReadColorIndexes(colorIndexesFilePath string) []stardata.ColorIndex {
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

	colorIndexes := createColorIndexesSlice(data)
	return colorIndexes
}

func createStarsSlice(data [][]string) []stardata.StarData {
	indexColumn, nameColumn, uColumn, bColumn, vColumn, sptypeColumn := -1, -1, -1, -1, -1, -1
	for j, field := range data[0] {
		clearField := strings.TrimSpace(field)
		switch clearField {
		case "#":
			indexColumn = j
		case "identifier":
			nameColumn = j
		case "Mag U":
			uColumn = j
		case "Mag B":
			bColumn = j
		case "Mag V":
			vColumn = j
		case "spec. type":
			sptypeColumn = j
		}
	}

	starsList := make([]stardata.StarData, 0, len(data))
	for i, line := range data {
		if i == 0 {
			continue
		}

		var index int
		var name, sptype string
		var u, b, v float64

		for j, field := range line {
			clearField := strings.TrimSpace(field)
			switch j {
			case indexColumn:
				num, err := strconv.Atoi(clearField)
				if err != nil {
					break
				}
				index = num

			case nameColumn:
				name = clearField

			case uColumn:
				value, err := strconv.ParseFloat(clearField, 64)
				if err != nil {
					break
				}
				u = value

			case bColumn:
				value, err := strconv.ParseFloat(clearField, 64)
				if err != nil {
					break
				}
				b = value

			case vColumn:
				value, err := strconv.ParseFloat(clearField, 64)
				if err != nil {
					break
				}
				v = value

			case sptypeColumn:
				sptype = clearField
			}
		}

		if u == 0 || b == 0 || v == 0 {
			continue
		}

		mag := stardata.NewMagnitude(u, b, v)
		sd := stardata.New(index, name, sptype, mag)

		starsList = append(starsList, *sd)
	}

	return starsList
}

func createColorIndexesSlice(data [][]string) []stardata.ColorIndex {
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

	colorIndexes := make([]stardata.ColorIndex, 0, len(data))
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

		colorIndexes = append(colorIndexes, stardata.NewColorIndex(colorIndex[0], colorIndex[1]))
	}

	return colorIndexes
}

func ReadMagVToBV(filePath string) [][2]float64 {
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

func WriteSliceToFile(filePath string, data [][2]float64) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Can't create file: %v", err)
	}
	defer file.Close()

	for _, chunk := range data {
		_, err := fmt.Fprintf(file, "%.4f\t%.4f\n", chunk[0], chunk[1])
		if err != nil {
			log.Printf("Can't write chunk to file: %v\n", err)
		}
	}
}

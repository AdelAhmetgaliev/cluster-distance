package stardata

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func CreateStarsList(data [][]string) (starsList []StarData) {
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

		mag := NewMagnitude(u, b, v)
		sd := New(index, name, sptype, mag)

		starsList = append(starsList, *sd)
	}

	return starsList
}

func WriteColorIndexesToFile(filePath string, data []StarData) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Can't create file: %v", err)
	}
	defer file.Close()

	for _, chunk := range data {
		_, err := fmt.Fprintf(file, "%.4f\t%.4f\n", chunk.CI.BV, chunk.CI.UB)
		if err != nil {
			log.Printf("Can't write chunk to file: %v\n", err)
		}
	}
}

func WriteMagVToBVToFile(filePath string, data []StarData) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Can't create file: %v", err)
	}
	defer file.Close()

	for _, chunk := range data {
		_, err := fmt.Fprintf(file, "%.4f\t%.4f\n", chunk.CI.BV, chunk.Mag.V)
		if err != nil {
			log.Printf("Can't write chunk to file: %v\n", err)
		}
	}
}

func AverageColorIndexes(data []StarData) ColorIndex {
	averageBV, averageUB := 0.0, 0.0
	for _, sd := range data {
		averageBV += sd.CI.BV
		averageUB += sd.CI.UB
	}

	averageBV /= float64(len(data))
	averageUB /= float64(len(data))
	averageCI := NewColorIndex(averageBV, averageUB)

	return averageCI
}

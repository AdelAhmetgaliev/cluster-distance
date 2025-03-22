package stardata

import (
	"fmt"
	"log"
	"os"
)

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
	aBV, aUB := 0.0, 0.0
	for _, sd := range data {
		aBV += sd.CI.BV
		aUB += sd.CI.UB
	}

	aBV /= float64(len(data))
	aUB /= float64(len(data))
	aCI := NewColorIndex(aBV, aUB)

	return aCI
}

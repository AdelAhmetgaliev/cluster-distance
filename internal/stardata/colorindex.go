package stardata

import (
	"fmt"
	"log"
	"os"
)

type ColorIndex struct {
	BV float64
	UB float64
}

func NewColorIndex(bv float64, ub float64) ColorIndex {
	return ColorIndex{bv, ub}
}

func (ci *ColorIndex) WriteToFile(filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Can't create file: %v", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%.4f\t%.4f\n", ci.BV, ci.UB)
	if err != nil {
		log.Printf("Can't write chunk to file: %v\n", err)
	}
}

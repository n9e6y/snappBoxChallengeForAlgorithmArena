package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"SBCFAA/internal/fare"
)

// WriteCSV writes the fare estimates to a CSV file
func WriteCSV(filename string, estimates []fare.FareEstimate) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"id_delivery", "fare_estimate"}); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write data
	for _, estimate := range estimates {
		row := []string{
			strconv.FormatInt(estimate.DeliveryID, 10),
			strconv.FormatFloat(estimate.Fare, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing row: %v", err)
		}
	}

	return nil
}

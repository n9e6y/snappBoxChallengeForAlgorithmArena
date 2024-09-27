package output

import (
	"encoding/csv"
	"io"
	"strconv"

	"SBCFAA/internal/models"
)

// WriteCSV writes fare estimates to a CSV file
func WriteCSV(writer io.Writer, estimates []models.FareEstimate) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	if err := csvWriter.Write([]string{"id_delivery", "fare_estimate"}); err != nil {
		return err
	}

	// Write data
	for _, estimate := range estimates {
		record := []string{
			strconv.FormatInt(estimate.DeliveryID, 10),
			strconv.FormatFloat(estimate.Fare, 'f', 2, 64),
		}
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}

	return nil
}

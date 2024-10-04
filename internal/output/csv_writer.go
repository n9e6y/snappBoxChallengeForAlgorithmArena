package output

import (
	"SBCFAA/internal/models"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"
)

const bufferSize = 1000 //change buffer size

func WriteCSV(filename string, estimates <-chan models.FareEstimate) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"id_delivery", "fare_estimate"}); err != nil { // Write header
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		buffer := make([][]string, 0, bufferSize)

		for estimate := range estimates {
			buffer = append(buffer, []string{
				strconv.FormatInt(estimate.DeliveryID, 10),
				strconv.FormatFloat(estimate.Fare, 'f', 2, 64),
			})

			if len(buffer) >= bufferSize {
				if err := writer.WriteAll(buffer); err != nil {
					// Log the error, but continue processing
					log.Printf("Error writing to CSV: %v", err)
				}
				buffer = buffer[:0] // Clear the buffer
			}
		}

		if len(buffer) > 0 { // Write any remaining records
			if err := writer.WriteAll(buffer); err != nil {
				log.Printf("Error writing final buffer to CSV: %v", err)
			}
		}
	}()

	wg.Wait()
	return nil
}

package output

import (
	"SBCFAA/internal/models"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"
)

const bufferSize = 1000

func WriteCSV(filename string, estimates <-chan models.FareEstimate) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"id_delivery", "fare_estimate"}); err != nil {
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

		// Write any remaining records
		if len(buffer) > 0 {
			if err := writer.WriteAll(buffer); err != nil {
				log.Printf("Error writing final buffer to CSV: %v", err)
			}
		}
	}()

	wg.Wait()
	return nil
}

//func WriteCSV(filename string, estimates <-chan []models.FareEstimate) error {
//	file, err := os.Create(filename)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	writer := csv.NewWriter(file)
//	defer writer.Flush()
//
//	// Write header
//	if err := writer.Write([]string{"id_delivery", "fare_estimate"}); err != nil {
//		return err
//	}
//
//	var wg sync.WaitGroup
//	wg.Add(1)
//
//	go func() {
//		defer wg.Done()
//		buffer := make([][]string, 0, bufferSize)
//
//		for chunk := range estimates {
//			for _, estimate := range chunk {
//				buffer = append(buffer, []string{
//					strconv.FormatInt(estimate.DeliveryID, 10),
//					strconv.FormatFloat(estimate.Fare, 'f', 2, 64),
//				})
//
//				if len(buffer) >= bufferSize {
//					if err := writer.WriteAll(buffer); err != nil {
//						// Log the error, but continue processing
//						log.Printf("Error writing to CSV: %v", err)
//					}
//					buffer = buffer[:0] // Clear the buffer
//				}
//			}
//		}
//
//		// Write any remaining records
//		if len(buffer) > 0 {
//			if err := writer.WriteAll(buffer); err != nil {
//				log.Printf("Error writing final buffer to CSV: %v", err)
//			}
//		}
//	}()
//
//	wg.Wait()
//	return nil
//}

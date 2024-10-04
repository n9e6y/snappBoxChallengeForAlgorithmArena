package ingestion

import (
	"SBCFAA/pkg/utils"
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"SBCFAA/internal/models"
)

func ReadAndFilterCSV(filename string) (<-chan []models.DeliveryPoint, <-chan error) {
	pointsChan := make(chan []models.DeliveryPoint, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(pointsChan)
		defer close(errChan)

		file, err := os.Open(filename)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()

		reader := csv.NewReader(bufio.NewReader(file))
		// Skip header
		if _, err := reader.Read(); err != nil {
			errChan <- err
			return
		}

		var currentDelivery []models.DeliveryPoint
		var currentID int64 = -1

		for {
			record, err := reader.Read()
			if err == io.EOF {
				if len(currentDelivery) > 0 {
					pointsChan <- currentDelivery
				}
				break
			}
			if err != nil {
				errChan <- err
				return
			}

			point, err := parseDeliveryPoint(record)
			if err != nil {
				errChan <- fmt.Errorf("error parsing record: %v", err)
				continue
			}

			if point.ID != currentID {
				if len(currentDelivery) > 0 {
					pointsChan <- currentDelivery
				}
				currentDelivery = []models.DeliveryPoint{point}
				currentID = point.ID
			} else {
				if len(currentDelivery) > 0 {
					prevPoint := currentDelivery[len(currentDelivery)-1]
					speed := calculateSpeed(prevPoint, point)
					if speed <= 100 { // 100 km/h filter
						currentDelivery = append(currentDelivery, point)
					}
				} else {
					currentDelivery = append(currentDelivery, point)
				}
			}
		}
	}()

	return pointsChan, errChan
}

func parseDeliveryPoint(record []string) (models.DeliveryPoint, error) {
	if len(record) != 4 {
		return models.DeliveryPoint{}, fmt.Errorf("invalid record length")
	}

	id, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return models.DeliveryPoint{}, err
	}

	lat, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return models.DeliveryPoint{}, err
	}

	lng, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return models.DeliveryPoint{}, err
	}

	timestamp, err := strconv.ParseInt(record[3], 10, 64)
	if err != nil {
		return models.DeliveryPoint{}, err
	}

	return models.DeliveryPoint{
		ID:        id,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Unix(timestamp, 0),
	}, nil
}

func calculateSpeed(p1, p2 models.DeliveryPoint) float64 {
	return utils.CalculateSpeed(p1, p2)
}

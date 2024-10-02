package ingestion

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"SBCFAA/pkg/utils"
)

type DeliveryPoint struct {
	ID        int64
	Lat, Lng  float64
	Timestamp int64
}

func ReadAndFilterCSV(filename string) ([]DeliveryPoint, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 4 // We expect 4 fields per record

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	var filteredData []DeliveryPoint
	var prevPoint *DeliveryPoint

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		point, err := parseRecord(record)
		if err != nil {
			return nil, fmt.Errorf("error parsing record: %v", err)
		}

		if prevPoint != nil && prevPoint.ID == point.ID {
			speed := calculateSpeed(*prevPoint, point)
			if speed <= 100 { // 100 km/h filter
				filteredData = append(filteredData, point)
			}
		} else {
			filteredData = append(filteredData, point)
		}

		prevPoint = &point
	}

	return filteredData, nil
}

func parseRecord(record []string) (DeliveryPoint, error) {
	id, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return DeliveryPoint{}, fmt.Errorf("invalid id: %v", err)
	}

	lat, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return DeliveryPoint{}, fmt.Errorf("invalid latitude: %v", err)
	}

	lng, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return DeliveryPoint{}, fmt.Errorf("invalid longitude: %v", err)
	}

	timestamp, err := strconv.ParseInt(record[3], 10, 64)
	if err != nil {
		return DeliveryPoint{}, fmt.Errorf("invalid timestamp: %v", err)
	}

	return DeliveryPoint{ID: id, Lat: lat, Lng: lng, Timestamp: timestamp}, nil
}

func calculateSpeed(p1, p2 DeliveryPoint) float64 {
	distance := utils.HaversineDistance(p1.Lat, p1.Lng, p2.Lat, p2.Lng)
	duration := time.Duration(p2.Timestamp-p1.Timestamp) * time.Second
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return distance / hours // km/h
}

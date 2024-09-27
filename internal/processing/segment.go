package processing

import (
	"SBCFAA/internal/models"
	"SBCFAA/pkg/utils"
	"time"
)

type Segment struct {
	Start    models.DeliveryPoint
	End      models.DeliveryPoint
	Distance float64
	Duration time.Duration
	Speed    float64
}

func ProcessSegments(points []models.DeliveryPoint) []Segment {
	segments := make([]Segment, 0, len(points)-1)

	for i := 1; i < len(points); i++ {
		start := points[i-1]
		end := points[i]

		distance := utils.HaversineDistance(start.Latitude, start.Longitude, end.Latitude, end.Longitude)
		duration := end.Timestamp.Sub(start.Timestamp)
		speed := distance / duration.Hours()

		segment := Segment{
			Start:    start,
			End:      end,
			Distance: distance,
			Duration: duration,
			Speed:    speed,
		}

		segments = append(segments, segment)
	}

	return segments
}

package domain

import (
	"time"
)

type ArgusMetric struct {
	Name      string
	Value     float64
	Timestamp time.Time
	Labels    Labels
}

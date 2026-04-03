package collector

import (
	"github.com/noboaki/argus-agent/domain"
)

type Collector interface {
	Collect() ([]*domain.ArgusMetric, error)
}

package processor

import "github.com/noboaki/argus-agent/domain"

type Processor interface {
	Process(*domain.ArgusMetric)
}

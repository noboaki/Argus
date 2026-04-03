package collector

import (
	"fmt"
	"time"

	"github.com/noboaki/argus-agent/domain"
	"github.com/shirou/gopsutil/v4/net"
)

type NetworkCollector struct {
	Interface string
	prev      *net.IOCountersStat
	prevTime  time.Time
}

func (n *NetworkCollector) Collect() ([]*domain.ArgusMetric, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("network collect error: %w", err)
	}

	now := time.Now()

	var current net.IOCountersStat
	for _, s := range stats {
		if n.Interface == "" || s.Name == n.Interface {
			current.BytesSent += s.BytesSent
			current.BytesRecv += s.BytesRecv
			current.PacketsSent += s.PacketsSent
			current.PacketsRecv += s.PacketsRecv
			current.Errin += s.Errin
			current.Errout += s.Errout
			current.Dropin += s.Dropin
			current.Dropout += s.Dropout
		}
	}

	if n.prev == nil {
		n.prev = &current
		n.prevTime = now
		return nil, nil
	}

	elapsed := now.Sub(n.prevTime).Seconds()

	bytesSentPerSec := float64(current.BytesSent-n.prev.BytesSent) / elapsed
	bytesRecvPerSec := float64(current.BytesRecv-n.prev.BytesRecv) / elapsed

	n.prev = &current
	n.prevTime = now

	labels := domain.Labels{"interface": n.interfaceName()}

	return []*domain.ArgusMetric{
		domain.NewArgusMetric("network_bytes_sent_per_sec", bytesSentPerSec).WithLabels(labels),
		domain.NewArgusMetric("network_bytes_recv_per_sec", bytesRecvPerSec).WithLabels(labels),
		domain.NewArgusMetric("network_errors_in", float64(current.Errin)).WithLabels(labels),
		domain.NewArgusMetric("network_errors_out", float64(current.Errout)).WithLabels(labels),
		domain.NewArgusMetric("network_drop_in", float64(current.Dropin)).WithLabels(labels),
		domain.NewArgusMetric("network_drop_out", float64(current.Dropout)).WithLabels(labels),
	}, nil
}

func (n *NetworkCollector) interfaceName() string {
	if n.Interface == "" {
		return "all"
	}
	return n.Interface
}

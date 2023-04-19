package discovery

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

type metrics struct {
	peerCount  asyncint64.Gauge
	peerEvents syncint64.Counter
}

func (d *Discovery) WithMetrics() error {
	peerCount, err := meter.AsyncInt64().Gauge(
		"disc_peer_count",
		instrument.WithDescription("Number of peers currently in discovery"),
	)
	if err != nil {
		return err
	}

	peerEvents, err := meter.SyncInt64().Counter(
		"disc_peer_events",
		instrument.WithDescription("Number of peer events in discovery"),
	)
	if err != nil {
		return err
	}

	d.metrics = &metrics{
		peerCount:  peerCount,
		peerEvents: peerEvents,
	}

	err = meter.RegisterCallback(
		[]instrument.Asynchronous{peerCount},
		func(ctx context.Context) {
			d.observePeerCount(ctx)
		},
	)
	return err
}

func (d *Discovery) observePeerEvent(ctx context.Context, s status) {
	if d.metrics == nil {
		return
	}

	d.metrics.peerEvents.Add(ctx, 1, attribute.String("status", string(s)))
}

func (d *Discovery) observePeerCount(ctx context.Context) {
	if d.metrics == nil {
		return
	}

	d.metrics.peerCount.Observe(ctx, int64(d.set.Size()))
}

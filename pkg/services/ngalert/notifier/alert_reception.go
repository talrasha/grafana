package notifier

import (
	"context"
	"time"

	gokit_log "github.com/go-kit/kit/log"
	apimodels "github.com/grafana/alerting-api/pkg/api"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/alertmanager/provider"
	"github.com/prometheus/alertmanager/provider/mem"
	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"
)

// DefaultResolveTimeout is the default timeout use for resolving the alert
// if the end time for alert is not specified.
// TODO: should this be configurable?
const DefaultResolveTimeout = 30 * time.Minute

type AlertProvider struct {
	provider.Alerts
}

// NewAlertProvider returns AlertProvider that provides a method to translate
// Grafana alerts to Prometheus Alertmanager alerts before passing it ahead.
func NewAlertProvider(m types.Marker) (*AlertProvider, error) {
	alerts, err := mem.NewAlerts(context.Background(), m, 30*time.Minute, gokit_log.NewNopLogger())
	if err != nil {
		return nil, err
	}

	return &AlertProvider{Alerts: alerts}, nil
}

func (ap *AlertProvider) PutPostableAlert(postableAlerts apimodels.PostableAlerts) error {
	now := time.Now()
	alerts := make([]*types.Alert, 0, len(postableAlerts.PostableAlerts))
	for _, a := range postableAlerts.PostableAlerts {
		alerts = append(alerts, alertForDelivery(a, now))
	}
	return ap.Alerts.Put(alerts...)
}

func alertForDelivery(a models.PostableAlert, now time.Time) *types.Alert {
	lbls := model.LabelSet{}
	annotations := model.LabelSet{}
	for k, v := range a.Labels {
		lbls[model.LabelName(k)] = model.LabelValue(v)
	}
	for k, v := range a.Annotations {
		annotations[model.LabelName(k)] = model.LabelValue(v)
	}

	alert := &types.Alert{
		Alert: model.Alert{
			Labels:       lbls,
			Annotations:  annotations,
			StartsAt:     time.Time(a.StartsAt),
			EndsAt:       time.Time(a.EndsAt),
			GeneratorURL: a.GeneratorURL.String(),
		},
		UpdatedAt: now,
	}

	// Ensure StartsAt is set.
	if alert.StartsAt.IsZero() {
		if alert.EndsAt.IsZero() {
			alert.StartsAt = now
		} else {
			alert.StartsAt = alert.EndsAt
		}
	}
	// If no end time is defined, set a timeout after which an alert
	// is marked resolved if it is not updated.
	if alert.EndsAt.IsZero() {
		alert.Timeout = true
		alert.EndsAt = now.Add(DefaultResolveTimeout)
	}

	return alert
}

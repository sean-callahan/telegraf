package kube_state

import (
	"context"
	"time"

	"github.com/influxdata/telegraf"
	"k8s.io/api/apps/v1beta2"
)

var (
	daemonSetMeasurement       = "kube_daemonset"
	daemonSetStatusMeasurement = "kube_daemonset_status"
)

func registerDaemonSetCollector(ctx context.Context, acc telegraf.Accumulator, ks *KubenetesState) {
	list, err := ks.client.getDaemonSets(ctx)
	if err != nil {
		acc.AddError(err)
		return
	}
	for _, d := range list.Items {
		if err = ks.gatherDaemonSet(d, acc); err != nil {
			acc.AddError(err)
			return
		}
	}
}

func (ks *KubenetesState) gatherDaemonSet(d v1beta2.DaemonSet, acc telegraf.Accumulator) error {
	if d.CreationTimestamp.IsZero() {
		return nil
	} else if !ks.firstTimeGather &&
		ks.MaxDaemonSetAge.Duration < time.Now().Sub(d.CreationTimestamp.Time) {
		return ks.gatherDaemonSetStatus(d, acc)
	}

	creationTime := d.CreationTimestamp.Time
	fields := map[string]interface{}{
		"metadata_generation": d.ObjectMeta.Generation,
	}
	tags := map[string]string{
		"namespace": d.Namespace,
		"daemonset": d.Name,
	}
	for k, v := range d.Labels {
		tags["label_"+sanitizeLabelName(k)] = v
	}

	acc.AddFields(daemonSetMeasurement, fields, tags, creationTime)
	return nil
}

func (ks *KubenetesState) gatherDaemonSetStatus(d v1beta2.DaemonSet, acc telegraf.Accumulator) error {
	status := d.Status
	fields := map[string]interface{}{
		"current_number_scheduled": status.CurrentNumberScheduled,
		"desired_number_scheduled": status.DesiredNumberScheduled,
		"number_available":         status.NumberAvailable,
		"number_misscheduled":      status.NumberMisscheduled,
		"number_ready":             status.NumberReady,
		"number_unavailable":       status.NumberUnavailable,
		"updated_number_scheduled": status.UpdatedNumberScheduled,
	}
	tags := map[string]string{
		"namespace": d.Namespace,
		"daemonset": d.Name,
	}
	acc.AddFields(daemonSetStatusMeasurement, fields, tags)
	return nil
}

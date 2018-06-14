package kube_state

import (
	"context"
	"time"

	"github.com/influxdata/telegraf"
	"k8s.io/api/core/v1"
)

var (
	replicationControllerMeasurement       = "kube_replicationcontroller"
	replicationControllerStatusMeasurement = "kube_replicationcontroller_status"
)

func registerReplicationControllerCollector(ctx context.Context, acc telegraf.Accumulator, ks *KubenetesState) {
	list, err := ks.client.getReplicationControllers(ctx)
	if err != nil {
		acc.AddError(err)
		return
	}
	for _, d := range list.Items {
		if err = ks.gatherReplicationController(d, acc); err != nil {
			acc.AddError(err)
			return
		}
	}
}

func (ks *KubenetesState) gatherReplicationController(d v1.ReplicationController, acc telegraf.Accumulator) error {
	var creationTime time.Time
	if !d.CreationTimestamp.IsZero() {
		creationTime = d.CreationTimestamp.Time
	}
	fields := map[string]interface{}{
		"metadata_generation": d.ObjectMeta.Generation,
	}
	tags := map[string]string{
		"namespace":             d.Namespace,
		"replicationcontroller": d.Name,
	}
	if d.Spec.Replicas != nil {
		fields["spec_replicas"] = *d.Spec.Replicas

	}
	acc.AddFields(replicationControllerMeasurement, fields, tags, creationTime)
	return ks.gatherReplicationControllerStatus(d, acc)
}

func (ks *KubenetesState) gatherReplicationControllerStatus(d v1.ReplicationController, acc telegraf.Accumulator) error {
	status := d.Status
	fields := map[string]interface{}{
		"replicas":               status.Replicas,
		"fully_labeled_replicas": status.FullyLabeledReplicas,
		"ready_replicas":         status.ReadyReplicas,
		"available_replicas":     status.AvailableReplicas,
		"observed_generation":    status.ObservedGeneration,
	}
	tags := map[string]string{
		"namespace":             d.Namespace,
		"replicationcontroller": d.Name,
	}
	acc.AddFields(replicationControllerStatusMeasurement, fields, tags)
	return nil
}

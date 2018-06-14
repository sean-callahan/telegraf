package kube_state

import (
	"context"

	"github.com/influxdata/telegraf"
	"k8s.io/api/apps/v1beta1"
)

var (
	statefulSetMeasurement       = "kube_statefulset"
	statefulSetStatusMeasurement = "kube_statefulset_status"
)

func registerStatefulSetCollector(ctx context.Context, acc telegraf.Accumulator, ks *KubenetesState) {
	list, err := ks.client.getStatefulSets(ctx)
	if err != nil {
		acc.AddError(err)
		return
	}
	for _, s := range list.Items {
		if err = ks.gatherStatefulSet(s, acc); err != nil {
			acc.AddError(err)
			return
		}
	}
}

func (ks *KubenetesState) gatherStatefulSet(statefulSet v1beta1.StatefulSet, acc telegraf.Accumulator) error {
	if statefulSet.CreationTimestamp.IsZero() {
		return nil
	}
	fields := map[string]interface{}{
		"metadata_generation": statefulSet.ObjectMeta.Generation,
	}
	tags := map[string]string{
		"namespace":   statefulSet.Namespace,
		"statefulset": statefulSet.Name,
	}
	if statefulSet.Spec.Replicas != nil {
		fields["replicas"] = *statefulSet.Spec.Replicas
	}
	for k, v := range statefulSet.Labels {
		tags["label_"+sanitizeLabelName(k)] = v
	}
	acc.AddFields(statefulSetMeasurement, fields, tags, statefulSet.CreationTimestamp.Time)
	return ks.gatherStatefulSetStatus(statefulSet, acc)
}

func (ks *KubenetesState) gatherStatefulSetStatus(statefulSet v1beta1.StatefulSet, acc telegraf.Accumulator) error {
	status := statefulSet.Status
	fields := map[string]interface{}{
		"replicas":         status.Replicas,
		"replicas_current": status.CurrentReplicas,
		"replicas_ready":   status.ReadyReplicas,
		"replicas_updated": status.UpdatedReplicas,
	}
	if statefulSet.Status.ObservedGeneration != nil {
		fields["observed_generation"] = *statefulSet.Status.ObservedGeneration
	}
	tags := map[string]string{
		"namespace":   statefulSet.Namespace,
		"statefulset": statefulSet.Name,
	}
	acc.AddFields(statefulSetStatusMeasurement, fields, tags)
	return nil
}

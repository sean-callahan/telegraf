package kube_state

import (
	"context"
	"time"

	"github.com/influxdata/telegraf"
	"k8s.io/api/core/v1"
)

var resourceQuotaMeasurement = "kube_resourcequota"

func registerResourceQuotaCollector(ctx context.Context, acc telegraf.Accumulator, ks *KubenetesState) {
	list, err := ks.client.getResourceQuotas(ctx)
	if err != nil {
		acc.AddError(err)
		return
	}
	for _, rq := range list.Items {
		if err = ks.gatherResourceQuota(rq, acc); err != nil {
			acc.AddError(err)
			return
		}
	}
}

func (ks *KubenetesState) gatherResourceQuota(rq v1.ResourceQuota, acc telegraf.Accumulator) error {
	var creationTime time.Time
	if !rq.CreationTimestamp.IsZero() {
		creationTime = rq.CreationTimestamp.Time
	}
	tags := map[string]string{
		"namespace":     rq.Namespace,
		"resourcequota": rq.Name,
	}

	for res, qty := range rq.Status.Hard {
		fields := map[string]interface{}{
			"gauge": qty.MilliValue() / 1000,
		}
		tags["resource"] = string(res)
		tags["type"] = "hard"
		acc.AddFields(resourceQuotaMeasurement, fields, tags, creationTime)
	}
	for res, qty := range rq.Status.Used {
		fields := map[string]interface{}{
			"gauge": qty.MilliValue() / 1000,
		}
		tags["resource"] = string(res)
		tags["type"] = "used"
		acc.AddFields(resourceQuotaMeasurement, fields, tags, creationTime)
	}

	return nil
}

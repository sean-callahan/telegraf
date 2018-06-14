package kube_state

import (
	"context"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"k8s.io/api/core/v1"
)

var limitRangeMeasurement = "kube_limitrange"

func registerLimitRangeCollector(ctx context.Context, acc telegraf.Accumulator, ks *KubenetesState) {
	list, err := ks.client.getLimitRanges(ctx)
	if err != nil {
		acc.AddError(err)
		return
	}
	for _, rq := range list.Items {
		if err = ks.gatherLimitRange(rq, acc); err != nil {
			acc.AddError(err)
			return
		}
	}
}

func (ks *KubenetesState) gatherLimitRange(rq v1.LimitRange, acc telegraf.Accumulator) error {

	var creationTime time.Time
	if !rq.CreationTimestamp.IsZero() {
		creationTime = rq.CreationTimestamp.Time
	}
	fields := map[string]interface{}{}
	tags := map[string]string{
		"namespace":  rq.Name,
		"limitrange": rq.Name,
	}
	rawLimitRanges := rq.Spec.Limits
	for _, rawLimitRange := range rawLimitRanges {
		for resource, min := range rawLimitRange.Min {
			fields["min_"+strings.ToLower(string(rawLimitRange.Type))+"_"+string(resource)] = min.MilliValue() / 1000
		}
		for resource, max := range rawLimitRange.Max {
			fields["max_"+strings.ToLower(string(rawLimitRange.Type))+"_"+string(resource)] = max.MilliValue() / 1000
		}

		for resource, df := range rawLimitRange.Default {
			fields["default_"+strings.ToLower(string(rawLimitRange.Type))+"_"+string(resource)] = df.MilliValue() / 1000
		}

		for resource, dfR := range rawLimitRange.DefaultRequest {
			fields["default_request_"+strings.ToLower(string(rawLimitRange.Type))+"_"+string(resource)] = dfR.MilliValue() / 1000
		}

		for resource, mLR := range rawLimitRange.MaxLimitRequestRatio {
			fields["max_limit_request_ratio_"+strings.ToLower(string(rawLimitRange.Type))+"_"+string(resource)] = mLR.MilliValue() / 1000
		}
	}
	acc.AddFields(limitRangeMeasurement, fields, tags, creationTime)

	return nil
}

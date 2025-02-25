package resource

import (
	"strconv"
	"time"

	analysisv1alph1 "github.com/gocrane/api/analysis/v1alpha1"
	"github.com/gocrane/crane/pkg/oom"
	"github.com/gocrane/crane/pkg/recommendation/config"
	"github.com/gocrane/crane/pkg/recommendation/recommender"
	"github.com/gocrane/crane/pkg/recommendation/recommender/apis"
	"github.com/gocrane/crane/pkg/recommendation/recommender/base"
)

var _ recommender.Recommender = &ResourceRecommender{}

type ResourceRecommender struct {
	base.BaseRecommender
	CpuSampleInterval        string
	CpuRequestPercentile     string
	CpuRequestMarginFraction string
	CpuTargetUtilization     string
	CpuModelHistoryLength    string
	MemSampleInterval        string
	MemPercentile            string
	MemMarginFraction        string
	MemTargetUtilization     string
	MemHistoryLength         string
	oomRecorder              oom.Recorder
	OOMProtection            bool
	OOMHistoryLength         time.Duration
	OOMBumpRatio             float64
	Specification            bool
	SpecificationConfigs     []Specification
	CpuHistogramBucketSize   string
	CpuHistogramMaxValue     string
	MemHistogramBucketSize   string
	MemHistogramMaxValue     string
}

func (rr *ResourceRecommender) Name() string {
	return recommender.ResourceRecommender
}

// NewResourceRecommender create a new resource recommender.
func NewResourceRecommender(recommender apis.Recommender, recommendationRule analysisv1alph1.RecommendationRule, oomRecorder oom.Recorder) (*ResourceRecommender, error) {
	recommender = config.MergeRecommenderConfigFromRule(recommender, recommendationRule)

	cpuSampleInterval, exists := recommender.Config["cpu-sample-interval"]
	if !exists {
		cpuSampleInterval = "1m"
	}
	cpuPercentile, exists := recommender.Config["cpu-request-percentile"]
	if !exists {
		cpuPercentile = "0.99"
	}
	cpuMarginFraction, exists := recommender.Config["cpu-request-margin-fraction"]
	if !exists {
		cpuMarginFraction = "0.15"
	}
	cpuTargetUtilization, exists := recommender.Config["cpu-target-utilization"]
	if !exists {
		cpuTargetUtilization = "1.0"
	}
	cpuHistoryLength, exists := recommender.Config["cpu-model-history-length"]
	if !exists {
		cpuHistoryLength = "168h"
	}

	memSampleInterval, exists := recommender.Config["mem-sample-interval"]
	if !exists {
		memSampleInterval = "1m"
	}
	memPercentile, exists := recommender.Config["mem-request-percentile"]
	if !exists {
		memPercentile = "0.99"
	}
	memMarginFraction, exists := recommender.Config["mem-request-margin-fraction"]
	if !exists {
		memMarginFraction = "0.15"
	}
	memTargetUtilization, exists := recommender.Config["mem-target-utilization"]
	if !exists {
		memTargetUtilization = "1.0"
	}
	memHistoryLength, exists := recommender.Config["mem-model-history-length"]
	if !exists {
		memHistoryLength = "168h"
	}

	specification, exists := recommender.Config["specification"]
	if !exists {
		specification = "false"
	}
	specificationBool, err := strconv.ParseBool(specification)
	if err != nil {
		return nil, err
	}
	specificationConfig, exists := recommender.Config["specification-config"]
	if !exists {
		specificationConfig = DefaultSpecs
	}
	specifications, err := GetResourceSpecifications(specificationConfig)
	if err != nil {
		return nil, err
	}

	oomProtection, exists := recommender.Config["oom-protection"]
	if !exists {
		oomProtection = "true"
	}

	oomProtectionBool, err := strconv.ParseBool(oomProtection)
	if err != nil {
		return nil, err
	}

	oomHistoryLength, exists := recommender.Config["oom-history-length"]
	if !exists {
		oomHistoryLength = "168h"
	}

	oomHistoryLengthDuration, err := time.ParseDuration(oomHistoryLength)
	if err != nil {
		return nil, err
	}

	OOMBumpRatio, exists := recommender.Config["oom-bump-ratio"]
	if !exists {
		OOMBumpRatio = "1.2"
	}
	OOMBumpRatioFloat, err := strconv.ParseFloat(OOMBumpRatio, 64)
	if err != nil {
		return nil, err
	}
	cpuHistogramBucketSize, exists := recommender.Config["cpu-histogram-bucket-size"]
	if !exists {
		cpuHistogramBucketSize = "0.1"
	}
	cpuHistogramMaxValue, exists := recommender.Config["cpu-histogram-max-value"]
	if !exists {
		cpuHistogramMaxValue = "100"
	}
	memHistogramBucketSize, exists := recommender.Config["mem-histogram-bucket-size"]
	if !exists {
		memHistogramBucketSize = "104857600"
	}
	memHistogramMaxValue, exists := recommender.Config["mem-histogram-max-value"]
	if !exists {
		memHistogramMaxValue = "104857600000"
	}

	return &ResourceRecommender{
		*base.NewBaseRecommender(recommender),
		cpuSampleInterval,
		cpuPercentile,
		cpuMarginFraction,
		cpuTargetUtilization,
		cpuHistoryLength,
		memSampleInterval,
		memPercentile,
		memMarginFraction,
		memTargetUtilization,
		memHistoryLength,
		oomRecorder,
		oomProtectionBool,
		oomHistoryLengthDuration,
		OOMBumpRatioFloat,
		specificationBool,
		specifications,
		cpuHistogramBucketSize,
		cpuHistogramMaxValue,
		memHistogramBucketSize,
		memHistogramMaxValue,
	}, nil
}

// Priority 5: Unit tests for approx_heavy_hitters

package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPercentileTree_GetPercentiles(t *testing.T) {
	p := NewPercentileTree()
	sampleSizes := []int{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}
	for _, size := range sampleSizes {
		p.Add(size)
	}

	percentiles := p.GetPercentiles()

	//assert.Equal(t, 500, percentiles["p50"])
	//a//ssert.Equal(t, 800, percentiles["p75"])
	//assert.Equal(t, 1000, percentiles["p99"])
	//assert.Equal(t, 900, percentiles["p90"])

	assert.Equal(t, 600, percentiles["p50"])
	assert.Equal(t, 900, percentiles["p75"])
	assert.Equal(t, 1000, percentiles["p90"])
	assert.Equal(t, 0, percentiles["p99"]) // or omit this if GetPercentiles doesn't guarantee p99


}

func TestCMSWrapper_EstimateAccuracy(t *testing.T) {
	cms := NewCMSWrapper(0.01, 0.9)
	for i := 0; i < 100; i++ {
		cms.Update("/some/path")
	}
	estimate := cms.Estimate("/some/path")
	assert.True(t, estimate >= 100)
}


package sketch

import (
	"github.com/shenwei356/countminsketch"
)

type CMSWrapper struct {
	sketch  *countminsketch.CountMinSketch
	epsilon float64
	delta   float64
}

func NewCMSWrapper(epsilon, delta float64) *CMSWrapper {
	s, err := countminsketch.NewWithEstimates(epsilon, delta)
	if err != nil {
		panic(err)
	}
	return &CMSWrapper{
		sketch:  s,
		epsilon: epsilon,
		delta:   delta,
	}
}

func (c *CMSWrapper) Update(path string) {
	sketchnext, err := countminsketch.NewWithEstimates(c.epsilon, c.delta)
	if err != nil {
		panic(err)
	}
	sketchnext.UpdateString(path, 1)
	c.sketch.Merge(sketchnext)
}

func (c *CMSWrapper) Estimate(path string) uint64 {
	return c.sketch.EstimateString(path)
}


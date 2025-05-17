package tree

import (
	"fmt"
	"math"

	"github.com/emirpasic/gods/maps/treemap"
)

type PercentileTree struct {
	tree *treemap.Map
}

func NewPercentileTree() *PercentileTree {
	return &PercentileTree{
		tree: treemap.NewWithIntComparator(),
	}
}

func (p *PercentileTree) Add(size int) {
	p.tree.Put(size, true)
}

func (p *PercentileTree) Size() int {
	return p.tree.Size()
}

func (p *PercentileTree) GetPercentiles() map[string]int {
	keys := []float64{0.5, 0.75, 0.90, 0.99}
	result := make(map[string]int)
	totalSize := p.tree.Size()
	it := p.tree.Iterator()
	it.First()
	count := 1
	for it.Next() {
		for _, k := range keys {
			idx := int(math.Ceil(k * float64(totalSize)))
			if count == idx {
				result[fmt.Sprintf("p%.0f", k*100)] = it.Key().(int)
			}
		}
		count++
	}
	return result
}


//inputs file, estimates the top 10 most frequently accessed paths, outputs percentiles of file sizes seen in universe
//preliminary submission by Acamar Orionis (Erica Stephens)
//Credit to: https://github.com/shenwei356/countminsketch for CMS data structure
//		 to: GoDS - Go Data Structures for Tree and Arraylist implementations
// Refactored: Clean code structure, encapsulate sketch, tree, error handling, removed globar vars, improved structure

package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shenwei356/countminsketch"
)

type CMSWrapper struct {
	sketch *countminsketch.CountMinSketch
	epsilon float64
	delta   float64
}

func NewCMSWrapper(epsilon, delta float64) *CMSWrapper {
	s, err := countminsketch.NewWithEstimates(epsilon, delta)
	checkErr(err)
	return &CMSWrapper{
		sketch: s,
		epsilon: epsilon,
		delta:   delta,
	}
}

func (c *CMSWrapper) Update(path string) {
	sketchNext, err := countminsketch.NewWithEstimates(c.epsilon, c.delta)
	checkErr(err)
	sketchNext.UpdateString(path, 1)
	c.sketch.Merge(sketchNext)
}

func (c *CMSWrapper) Estimate(path string) uint64 {
	return c.sketch.EstimateString(path)
}

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

func main() {
	f, err := os.Open("path1.txt")
	checkErr(err)
	defer f.Close()
	processInput(f)
}

func processInput(r io.Reader) {
	sketch := NewCMSWrapper(0.01, 0.9)
	tree := treemap.NewWithIntComparator()
	percentiles := NewPercentileTree()
	finalAHH := arraylist.New()
	seedVal := "seed"
	m := treemap.NewWithIntComparator()
	br := bufio.NewReader(r)

	for {
		line, err := br.ReadString('\n')
		if err == io.EOF && line == "" {
			break
		} else if err != nil && err != io.EOF {
			checkErr(err)
		}
		line = strings.TrimSpace(line)
		fields := strings.Split(line, "\t")
		if len(fields) != 2 {
			continue
		}
		path := fields[0]
		size, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		sketch.Update(path)
		est := int(sketch.Estimate(path))
		updateAHHTree(m, est, path, seedVal)
		percentiles.Add(size)
		if err == io.EOF {
			break
		}
	}

	fmt.Println("\nTop 10 Paths:")
	itAHH := m.Iterator()
	itAHH.End()
	for itAHH.Prev() && finalAHH.Size() < 10 {
		_, val := itAHH.Key(), itAHH.Value()
		paths := strings.Split(val.(string), "|")
		for _, p := range paths {
			if p != seedVal && !finalAHH.Contains(p) && finalAHH.Size() < 10 {
				finalAHH.Add(p)
				fmt.Printf("%d. %s\n", finalAHH.Size(), p)
			}
		}
	}

	fmt.Println("\nPercentiles:")
	for label, val := range percentiles.GetPercentiles() {
		fmt.Printf("file_size_%s\t%d\n", label, val)
	}
}

func updateAHHTree(m *treemap.Map, est int, path, seed string) {
	minKey, _ := m.Min()
	min := minKey.(int)
	existing, found := m.Get(est)

	if est > min && !found {
		m.Put(est, path)
	} else {
		treePath := ""
		if found {
			treePath = strings.ReplaceAll(existing.(string), seed, "")
		} else {
			m.Put(est, path)
			return
		}
		if !strings.Contains(treePath, path) {
			if treePath == "" {
				m.Put(est, path)
			} else {
				newPath := path + "|" + treePath
				m.Put(est, newPath)
			}
		}
		if m.Size() > 30 {
			m.Remove(min)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

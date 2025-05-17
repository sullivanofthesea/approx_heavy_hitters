//inputs file, estimates the top 10 most frequently accessed paths, outputs percentiles of file sizes seen in universe
//preliminary submission by Acamar Orionis (Erica Stephens)
//Credit to: https://github.com/shenwei356/countminsketch for CMS data structure
//		 to: GoDS - Go Data Structures for Tree and Arraylist implementations
// Refactored: Clean code structure, encapsulate sketch, tree, error handling, removed globar vars, improved structure
// Refactored: Add config file instead of hard coding (boo) - Load epsilon/delta from config.json
// Feature added: Simple batching and mem control (can be enhanced later)
// Feature added: CLI flags for input file, batchSize, varepsilon, delta

package main

import (
	"bufio"
	"encoding/json"
	"flag"
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

type Config struct {
	Varepsilon float64 `json:"varepsilon"`
	Delta      float64 `json:"delta"`
}

type CMSWrapper struct {
	sketch  *countminsketch.CountMinSketch
	epsilon float64
	delta   float64
}

func NewCMSWrapper(epsilon, delta float64) *CMSWrapper {
	s, err := countminsketch.NewWithEstimates(epsilon, delta)
	checkErr(err)
	return &CMSWrapper{
		sketch:  s,
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

func loadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		return Config{Varepsilon: 0.01, Delta: 0.9} // fallback
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return Config{Varepsilon: 0.01, Delta: 0.9}
	}
	return config
}

func main() {
	inputFile := flag.String("input", "path1.txt", "Path to input file")
	varepsilon := flag.Float64("varepsilon", 0.01, "Error bound for CMS")
	delta := flag.Float64("delta", 0.9, "Confidence level for CMS")
	batchSize := flag.Int("batchSize", 1000, "Number of lines per batch")
	flag.Parse()

	config := loadConfig("config.json")

	// CLI flags override config.json if provided
	if *varepsilon == 0 {
		*varepsilon = config.Varepsilon
	}
	if *delta == 0 {
		*delta = config.Delta
	}

	f, err := os.Open(*inputFile)
	checkErr(err)
	defer f.Close()

	processInput(f, *varepsilon, *delta, *batchSize)
}

func processInput(r io.Reader, varepsilon float64, delta float64, batchSize int) {
	lineCount := 0
	sketch := NewCMSWrapper(varepsilon, delta)
	percentiles := NewPercentileTree()
	finalAHH := arraylist.New()
	seedVal := "seed"
	m := treemap.NewWithIntComparator()
	br := bufio.NewScanner(r)

	for br.Scan() {
		line := strings.TrimSpace(br.Text())
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
		lineCount++

		if lineCount%batchSize == 0 {
			printBatchSummary(m, percentiles, finalAHH, seedVal)
			sketch = NewCMSWrapper(varepsilon, delta)
			percentiles = NewPercentileTree()
			m = treemap.NewWithIntComparator()
			finalAHH = arraylist.New()
		}
	}

	if lineCount%batchSize != 0 {
		printBatchSummary(m, percentiles, finalAHH, seedVal)
	}
}

func printBatchSummary(m *treemap.Map, percentiles *PercentileTree, finalAHH *arraylist.List, seedVal string) {
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
	if minKey == nil {
		m.Put(est, path)
		return
	}
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

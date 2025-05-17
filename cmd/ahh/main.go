//inputs file, estimates the top 10 most frequently accessed paths, outputs percentiles of file sizes seen in universe
//preliminary submission by Acamar Orionis (Erica Stephens)
//Credit to: https://github.com/shenwei356/countminsketch for CMS data structure
//		 to: GoDS - Go Data Structures for Tree and Arraylist implementations
// Refactored: Clean code structure, encapsulate sketch, tree, error handling, removed globar vars, improved structure
// Refactored: Add config file instead of hard coding (boo) - Load epsilon/delta from config.json
// Feature added: Simple batching and mem control (can be enhanced later)
// Feature added: CLI flags for input file, batchSize, varepsilon, delta
// Restructured: organized directory
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/estephensltu/approx_heavy_hitters/internal/ahh"
)

type Config struct {
	Varepsilon float64 `json:"varepsilon"`
	Delta      float64 `json:"delta"`
}

func main() {
	inputFile := flag.String("input", "path1.txt", "Path to input file")
	varepsilon := flag.Float64("varepsilon", 0.01, "Error bound for CMS")
	delta := flag.Float64("delta", 0.9, "Confidence level for CMS")
	batchSize := flag.Int("batchSize", 1000, "Number of lines per batch")
	flag.Parse()

	config := loadConfig("config/config.json")

	if *varepsilon == 0 {
		*varepsilon = config.Varepsilon
	}
	if *delta == 0 {
		*delta = config.Delta
	}

	file, err := os.Open(*inputFile)
	checkErr(err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	ahh.ProcessInput(lines, *varepsilon, *delta, *batchSize)
}

func loadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		return Config{Varepsilon: 0.01, Delta: 0.9}
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{Varepsilon: 0.01, Delta: 0.9}
	}
	return config
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}


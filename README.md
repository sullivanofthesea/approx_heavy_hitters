# approx_heavy_hitters

## Overview

`approx_heavy_hitters` is a command-line utility written in Go that estimates the most frequently accessed file paths (approximate heavy hitters) and computes statistical percentiles for file sizes within a stream of tab-delimited path–size records. It is designed for use with large input files or unbounded data streams where memory efficiency is critical and exact frequency tracking is impractical.

The tool uses a Count-Min Sketch (CMS) to estimate the frequency of each unique file path, trading off a small bounded error for constant memory usage. It maintains a TreeMap to track the top-k file paths with the highest estimated access frequencies and applies collision handling using chained path strings. For size analysis, a second TreeMap is used to store and sort file size values, allowing efficient calculation of approximate p50, p75, p90, and p99 percentiles.

This implementation is single-pass, operates entirely in memory, and includes placeholder functionality for batching, pruning, and socket-based output extensions. The current config processes input from a local file named `path1.txt` and prints output to the console. Code is structured with a focus on compact data processing, though TODO sections show areas for future encapsulation, configurability, modular design, and testing coverage.

## Requirements

- Go 1.18 or later
- A tab-delimited input file named `path1.txt` in the same directory

## Input Format

Each line should contain a path and a file size in bytes, separated by a tab:

```
/some/path/file.txt	1234
/another/path.log	5678
```

## Usage

```bash
go run main.go
```

or

```bash
go build -o ahh
./ahh
```

## Output

The program prints:

- Top 10 most frequent paths (approximate)
- 50th, 75th, 90th, and 99th percentiles of observed file sizes

## Notes

- Paths are estimated using a Count-Min Sketch
- Percentiles are derived from a TreeMap of file sizes

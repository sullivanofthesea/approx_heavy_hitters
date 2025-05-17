# approx_heavy_hitters

## Overview

`approx_heavy_hitters` is a command line utility written in Go that estimates the most frequently accessed file paths (approximate heavy hitters) and computes statistical percentiles for file sizes within a stream of tab delimited path size records.

It is designed for use with large input files or unbounded data streams where memory efficiency is critical and exact frequency tracking is impractical.

The tool uses a Count-Min Sketch (CMS) to estimate frequency of each unique file path, trading a small bounded error for constant memory usage. It maintains a TreeMap to track the top-k file paths with the highest estimated access frequencies and a second TreeMap to track file sizes for efficient percentile calculations (p50, p75, p90, p99).

This implementation is single pass, batch aware, operates entirely in memory, and is configurable via CLI flags or a JSON config file.

---

## Why Use This

In systems that handle large volumes of file access logs such as distributed storage systems, content delivery networks, or compute clusters it's often valuable to understand which files are accessed most frequently and how their sizes are distributed.

`approx_heavy_hitters` provides a memory efficient, approximate solution using Count Min Sketch to identify hot paths and track file size distributions without storing the entire dataset. This helps with:

- **Cache optimization**
- **Storage tiering**
- **Load balancing**
- **Shard rebalancing**
- **Monitoring and alerting pipelines**

Because it is single-pass and uses constant memory relative to unique keys, it’s suitable for log processing pipelines and realtime or large scale analytics.

---

## Requirements

- Go 1.18 or later
- A tab delimited input file (`path1.txt` by default)
- Optional: `config/config.json` for sketch parameters

---

## Input Format

Each line should contain a path and file size in bytes, separated by a tab:

```
/some/path/file.txt	1234
/another/path.log	5678
```

---

## Usage

### Run with Go:

```bash
go run ./cmd/ahh --input=path1.txt --batchSize=1000 --varepsilon=0.01 --delta=0.9
```

### Or build and run:

```bash
go build -o ahh ./cmd/ahh
./ahh --input=path1.txt
```

### Flags:

| Flag         | Description                             | Default     |
|--------------|-----------------------------------------|-------------|
| `--input`     | Path to input file                      | `path1.txt` |
| `--batchSize` | Number of lines to process per batch    | `1000`      |
| `--varepsilon`| CMS error bound                         | `0.01`      |
| `--delta`     | CMS confidence level                    | `0.9`       |

*Flags override values in `config/config.json`.*

---

## Output

After each batch, the program prints:

- Top 10 most frequent file paths (approximate)
- File size percentiles: p50, p75, p90, p99

---

## Project Structure

```
approx_heavy_hitters/
├── cmd/ahh/                # CLI entry
│   └── main.go
├── config/                 # Config and static files
│   └── config.json
├── internal/               # Core logic
│   ├── ahh/                # Batch and heavy hitter logic
│   ├── sketch/             # CMS wrapper
│   └── tree/               # Percentile logic
├── test/                   # Unit tests
│   └── main_test.go
├── tools/                  # Dev/test data generation
│   └── generate_input.go
├── go.mod
├── go.sum
└── README.md
```

---

## Notes

- Uses GoDS (TreeMap and ArrayList)
- CMS logic from `github.com/shenwei356/countminsketch`
- Batched processing with memory reset between batches
- No dependencies beyond standard Go modules and external packages listed in `go.mod`
- Tested via `go test ./...`

---

## TODO / Next Steps

- [x] Add test coverage for `updateAHHTree` and core logic
- [x] Refactor into internal packages for clarity and reuse
- [x] CLI flags and JSON config loading
- [x] Batch processing logic
- [x] Reorganize project by purpose (src, config, tools, test)
- [ ] Add CI pipeline using GitHub Actions
- [ ] Support output in JSON or exportable format
- [ ] Add `--topN` flag to control number of paths shown
- [ ] Add Makefile or build script for packaging

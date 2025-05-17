# approx_heavy_hitters

## Overview

`approx_heavy_hitters` is a command-line utility written in Go that estimates the most frequently accessed file paths (approximate heavy hitters) and computes statistical percentiles for file sizes within a stream of tab-delimited path–size records. It is designed for use with large input files or unbounded data streams where memory efficiency is critical and exact frequency tracking is impractical.

The tool uses a Count-Min Sketch (CMS) to estimate the frequency of each unique file path, trading off a small bounded error for constant memory usage. It maintains a TreeMap to track the top-k file paths with the highest estimated access frequencies and applies collision handling using chained path strings. For size analysis, a second TreeMap is used to store and sort file size values, allowing efficient calculation of approximate p50, p75, p90, and p99 percentiles.

This implementation is single-pass, operates entirely in memory, and includes placeholder functionality for batching, pruning, and socket-based output extensions. The current config processes input from a local file named `path1.txt` and prints output to the console. Code is structured with a focus on compact data processing, though TODO sections show areas for future encapsulation, configurability, modular design, and testing coverage.

## Why Use This

In systems that handle large volumes of file access logs, such as distributed storage systems, content delivery networks, or server farms, it can be important to understand which files are accessed most frequently and how file sizes are distributed. Performing exact tracking is often memory-intensive and impractical in real time.

`approx_heavy_hitters` provides a memory-efficient, approximate solution using Count-Min Sketch to identify high-frequency paths without storing the entire dataset. The percentile calculations for file sizes help in assessing storage patterns, detecting outliers, and informing system-level optimizations such as caching, tiered storage, or resource allocation.

Because the implementation is single-pass and requires constant memory relative to the number of unique paths, it is suitable for integration into monitoring pipelines, batch processing jobs, or log pre-processing stages in systems where throughput and scalability are critical.

## Operational Use Cases

This tool can be integrated into larger systems where understanding file access patterns is relevant to infrastructure planning or system performance.

### Load Balancing

In distributed storage systems, compute clusters, or caching architectures, access patterns are often skewed toward a small number of frequently requested files. `approx_heavy_hitters` can support load balancing efforts in the following ways:

- **Identifying Hotspots**: Approximates the most frequently accessed paths, allowing operators to detect traffic concentration and adjust routing or resource allocation accordingly.
- **Cache Optimization**: High-frequency paths can be prioritized for replication or caching across nodes to reduce latency and prevent overloading individual systems.
- **Shard Rebalancing**: Frequency data can be used to inform the reassignment of data or compute workloads across shards or partitions to improve utilization.
- **Storage Tiering**: Frequently accessed items can be elevated to faster storage tiers, while less-accessed items are moved to lower-cost storage.
- **Streamlining Monitoring**: Percentile metrics and access estimates can be incorporated into real-time analytics dashboards or automation scripts for operational decision-making.

The memory-bounded and single-pass nature of the implementation makes it suitable for log processing pipelines, batch analytics jobs, or inline preprocessing for large-scale event streams.

## Requirements

- Go 1.18 or later
- A tab-delimited input file named `path1.txt` in the same directory
- Optional: `config.json` for sketch parameters


## Input Format

Each line should contain a path and a file size in bytes, separated by a tab:

```
/some/path/file.txt	1234
/another/path.log	5678
```

## Usage

```bash
go run main.go --input=path1.txt --batchSize=1000 --varepsilon=0.01 --delta=0.9
```

Flags override values from config.json if both are provided.


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

- Batching: Processes input in batches with optional memory reset
- Count-Min Sketch: Estimates frequency of access per path
- Percentiles: Derived from a TreeMap of file sizes
- Configurable via CLI flags or optional config.json
- Unit-tested with `go test`

## TODO

- [x] Add test coverage for `updateAHHTree`
- [x] Package CLI with user-friendly flags/help text
- [x] Add percentile test harness with `go test`
- [ ] Add CI pipeline using GitHub Actions
- [ ] Add support for `--topN` or output to JSON
- [ ] Add Makefile or installer for CLI distribution



# approx_heavy_hitters

This tool estimates the most frequently accessed file paths and calculates file size percentiles. It uses Count-Min Sketch and TreeMap-based structures.

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

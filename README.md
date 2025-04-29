# Customer CLI

A command-line application that reads a CSV file of customer emails, counts domain usage, and outputs sorted results based on the selected options.

## Features

- Reads large CSV files efficiently using channels and goroutines.
- Extracts domains from email addresses.
- Counts domain occurrences.
- Supports sorting by domain name or usage count (ascending or descending).
- Outputs results to terminal or file.

## Usage

```bash
customercli [flags] <inputfile>
```

### Arguments

| Argument    | Description                                |
|-------------|--------------------------------------------|
| inputfile   | Path to the CSV file (required)             |

### Flags

| Flag               | Description                                                                                  |
|--------------------|----------------------------------------------------------------------------------------------|
| `--sort=name`       | Sort by "name" (alphabetical) or "count" (usage count). Default is `name`.                   |
| `--order=asc`       | Sort order: `asc` (ascending) or `desc` (descending). Optional. Default is `desc` if not provided for `count`. |
| `--output=filename` | Output file path. Optional. If omitted, output is printed to terminal.                      |
| `-h, --help`        | Show help message.                                                                           |

## Examples

```bash
# Sort by domain name (alphabetically)
customercli customers.csv

# Sort by count descending (most common domains first)
customercli --sort=count customers.csv

# Sort by count ascending
customercli --sort=count --order=asc customers.csv

# Output to a file
customercli --output=output.txt customers.csv

# Sort by count ascending and write to file
customercli --sort=count --order=asc --output=output.txt customers.csv
```

## Project Structure

```
.
├── cmd/
│   ├── main.go           # CLI entry
│   └── main_test.go      # CLI tests
├── customerimporter/
│   ├── interview.go      # Core logic
│   └── interview_test.go # Unit tests
├── go.mod
└── README.md

```

## How to Build

To build the CLI application:

```bash
go build -o customercli ./cmd
```

This will generate an executable named `customercli`.

## How to Run

After building, you can run the program like this:

```bash
./customercli customers.csv
```

or with flags:

```bash
./customercli --sort=count --order=asc --output=result.txt customers.csv
```

If you prefer not to build and just want to run it directly:

```bash
go run ./cmd --sort=count customers.csv
```

## Developer Notes

- Written in Go, using only the standard library (no external packages).
- Efficient memory handling using channels and goroutines while reading CSV files.

---

Built with 💻 and Go!

---

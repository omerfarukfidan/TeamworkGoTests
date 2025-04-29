package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"teamworkgotests/customerimporter"
)

func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `Usage of %s:

Usage:
  customercli [flags] <inputfile>

Arguments:
  inputfile          Path to the CSV file.

Flags:
  --sort=name        Sort domains by "name" (alphabetical) or "count" (number of uses). Default is "name".
  --order=asc        (Optional) When sorting by count, you can specify "asc" (ascending) or "desc" (descending).
                     If --sort=count is used without --order, the default is descending.
  --output=filename  Write output to a file instead of printing to terminal.
  -h, --help         Show help message.

Examples:
  customercli customers.csv
  customercli --sort=count customers.csv
  customercli --sort=count --order=asc customers.csv
  customercli --output=result.txt customers.csv
  customercli --sort=count --order=asc --output=result.txt customers.csv

`, os.Args[0])
	}

	log.SetFlags(0)
}

func main() {
	// Parse CLI flags
	sortBy := flag.String("sort", "name", "Sort by 'name' or 'count'")
	output := flag.String("output", "", "Output file path (optional)")
	order := flag.String("order", "", "Sort order: asc or desc (optional)")

	flag.Parse()

	*sortBy = strings.ToLower(*sortBy)
	*order = strings.ToLower(*order)

	if *order == "" && *sortBy == "count" {
		*order = "desc"
	}

	logInfo("Starting customercli...")

	// Check for CSV input file
	if flag.NArg() < 1 {
		logError("No input file provided.")
		_, _ = fmt.Fprintln(os.Stderr, "Use -h or --help for usage information.")
		os.Exit(1)
	}
	inputFile := flag.Arg(0)

	logInfo("Reading emails from: %s", inputFile)

	// Read emails from CSV
	emails, err := customerimporter.ReadEmailsFromCSV(inputFile)
	if err != nil {
		logError("Failed to read CSV file: %v", err)
		os.Exit(1)
	}

	// Count domains
	domainCounts, invalidCount := customerimporter.CountDomains(emails)

	logInfo("Total unique domains found: %d", len(domainCounts))
	logInfo("Total invalid emails skipped: %d", invalidCount)

	// Sort domains
	sortedDomains := customerimporter.SortDomainCounts(domainCounts, *sortBy, *order)

	// Determine output destination
	var writer *os.File
	if *output == "" {
		logInfo("Writing output to terminal...")
		writer = os.Stdout
	} else {
		logInfo("Writing output to file: %s", *output)
		writer, err = os.Create(*output)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer func() {
			if err := writer.Close(); err != nil {
				logWarning("Error closing output file: %v", err)
			}
		}()

	}

	// Print domain counts
	err = customerimporter.PrintDomainCounts(sortedDomains, writer)
	if err != nil {
		log.Fatalf("Failed to print domain counts: %v", err)
	}

	logInfo("Output successfully written.")
}

func logInfo(format string, args ...interface{}) {
	log.Printf("\033[34m[INFO]\033[0m "+format, args...)
}

func logError(format string, args ...interface{}) {
	log.Printf("\033[31m[ERROR]\033[0m "+format, args...)
}

func logWarning(format string, args ...interface{}) {
	log.Printf("\033[33m[WARNING]\033[0m "+format, args...)
}

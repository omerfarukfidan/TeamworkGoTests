package customerimporter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

// DomainCount holds a domain name and its count.
type DomainCount struct {
	Domain string
	Count  int
}

// ReadEmailsFromCSV reads email addresses from a CSV file and sends them into a channel.
// It reads line by line to avoid high memory usage with large files.
func ReadEmailsFromCSV(filename string) (<-chan string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	out := make(chan string)

	// Start a goroutine to read the file asynchronously
	go func() {
		defer func() {
			if err := file.Close(); err != nil {
				logWarning("Error closing file: %v", err)
			}
		}()
		defer close(out)

		isHeader := true

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading CSV record: %v", err)
				continue // Skip bad lines, but don't crash
			}

			if isHeader {
				isHeader = false
				continue // Skip header line
			}

			if len(record) < 3 {
				continue // skip if not enough fields
			}

			email := record[2]
			out <- email
		}
	}()

	return out, nil
}

// ExtractDomain extracts the domain from an email address.
// Returns an error if the email address is invalid.
func ExtractDomain(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", errors.New("invalid email address format")
	}
	return parts[1], nil
}

// CountDomains reads emails from a channel, extracts domains, and counts their occurrences.
// Returns a map where the key is the domain and the value is the number of occurrences.
func CountDomains(emails <-chan string) (map[string]int, int) {
	domainCounts := make(map[string]int)
	invalidCount := 0

	for email := range emails {
		domain, err := ExtractDomain(email)
		if err != nil {
			logWarning("Skipping invalid email (%s): %v", email, err)
			invalidCount++
			continue
		}
		domainCounts[domain]++
	}

	return domainCounts, invalidCount
}

// SortDomainCounts sorts the domain counts according to the given sort criteria ("name" or "count").
// It returns a sorted slice of DomainCount.
func SortDomainCounts(domainCounts map[string]int, sortBy, order string) []DomainCount {
	var sorted []DomainCount
	for domain, count := range domainCounts {
		sorted = append(sorted, DomainCount{
			Domain: domain,
			Count:  count,
		})
	}

	normalizedSortBy := strings.ToLower(sortBy)
	normalizedOrder := strings.ToLower(order)

	if normalizedSortBy == "count" {
		sort.Slice(sorted, func(i, j int) bool {
			if normalizedOrder == "asc" {
				return sorted[i].Count < sorted[j].Count
			}
			return sorted[i].Count > sorted[j].Count
		})
	} else {
		sort.Slice(sorted, func(i, j int) bool {
			if normalizedOrder == "desc" {
				return sorted[i].Domain > sorted[j].Domain
			}
			return sorted[i].Domain < sorted[j].Domain
		})
	}

	return sorted
}

// PrintDomainCounts prints the domain counts to the given writer (stdout or file).
func PrintDomainCounts(domainCounts []DomainCount, writer io.Writer) error {
	for _, dc := range domainCounts {
		_, err := fmt.Fprintf(writer, "%s: %d\n", dc.Domain, dc.Count)
		if err != nil {
			return err
		}
	}
	return nil
}

func logWarning(format string, args ...interface{}) {
	log.Printf("\033[33m[WARNING]\033[0m "+format, args...)
}

package customerimporter

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		email       string
		wantDomain  string
		expectError bool
	}{
		{"user@example.com", "example.com", false},
		{"someone@domain.org", "domain.org", false},
		{"invalid-email.com", "", true},
		{"anotherinvalidemail@", "", true},
		{"@nodomain", "", true},
	}

	for _, tt := range tests {
		domain, err := ExtractDomain(tt.email)
		if (err != nil) != tt.expectError {
			t.Errorf("ExtractDomain(%q) error = %v, wantError = %v", tt.email, err, tt.expectError)
		}
		if domain != tt.wantDomain {
			t.Errorf("ExtractDomain(%q) = %q, want %q", tt.email, domain, tt.wantDomain)
		}
	}
}

func TestCountDomains(t *testing.T) {
	emails := make(chan string)

	go func() {
		defer close(emails)
		emails <- "user1@example.com"
		emails <- "user2@example.com"
		emails <- "user3@test.com"
		emails <- "invalid-email.com"
		emails <- "user4@test.com"
	}()

	expected := map[string]int{
		"example.com": 2,
		"test.com":    2,
	}

	result, _ := CountDomains(emails)

	if len(result) != len(expected) {
		t.Fatalf("expected %d domains, got %d", len(expected), len(result))
	}

	for domain, count := range expected {
		if result[domain] != count {
			t.Errorf("for domain %q, expected count %d, got %d", domain, count, result[domain])
		}
	}
}

func TestSortDomainCounts(t *testing.T) {
	domainCounts := map[string]int{
		"example.com": 5,
		"test.com":    8,
		"alpha.com":   3,
	}

	t.Run("sort by name asc", func(t *testing.T) {
		sorted := SortDomainCounts(domainCounts, "name", "asc")

		expectedOrder := []string{"alpha.com", "example.com", "test.com"}

		for i, domain := range expectedOrder {
			if sorted[i].Domain != domain {
				t.Errorf("expected domain at index %d to be %q, got %q", i, domain, sorted[i].Domain)
			}
		}
	})

	t.Run("sort by name desc", func(t *testing.T) {
		sorted := SortDomainCounts(domainCounts, "name", "desc")

		expectedOrder := []string{"test.com", "example.com", "alpha.com"}

		for i, domain := range expectedOrder {
			if sorted[i].Domain != domain {
				t.Errorf("expected domain at index %d to be %q, got %q", i, domain, sorted[i].Domain)
			}
		}
	})

	t.Run("sort by count desc", func(t *testing.T) {
		sorted := SortDomainCounts(domainCounts, "count", "desc")

		expectedOrder := []string{"test.com", "example.com", "alpha.com"}

		for i, domain := range expectedOrder {
			if sorted[i].Domain != domain {
				t.Errorf("expected domain at index %d to be %q, got %q", i, domain, sorted[i].Domain)
			}
		}
	})

	t.Run("sort by count asc", func(t *testing.T) {
		sorted := SortDomainCounts(domainCounts, "count", "asc")

		expectedOrder := []string{"alpha.com", "example.com", "test.com"}

		for i, domain := range expectedOrder {
			if sorted[i].Domain != domain {
				t.Errorf("expected domain at index %d to be %q, got %q", i, domain, sorted[i].Domain)
			}
		}
	})
}

func TestPrintDomainCounts(t *testing.T) {
	data := []DomainCount{
		{"example.com", 2},
		{"test.com", 1},
	}
	var buf bytes.Buffer
	err := PrintDomainCounts(data, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "example.com: 2") || !strings.Contains(output, "test.com: 1") {
		t.Errorf("Unexpected output: %s", output)
	}
}

func TestReadEmailsFromCSV_FileNotFound(t *testing.T) {
	_, err := ReadEmailsFromCSV("nonexistent.csv")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestReadEmailsFromCSV_InvalidFormat(t *testing.T) {
	filename := "broken.csv"
	content := "id,name,email\n1,Omer\n2,Ali,ali@test.com\n"
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	defer func() {
		if err := os.Remove(filename); err != nil {
			t.Logf("warning: failed to remove temp file %s: %v", filename, err)
		}
	}()

	ch, err := ReadEmailsFromCSV(filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count := 0
	for range ch {
		count++
	}

	if count != 1 {
		t.Errorf("expected 1 valid email, got %d", count)
	}
}

func TestReadEmailsFromCSV_Valid(t *testing.T) {
	filename := "valid.csv"
	content := "id,name,email\n1,Omer,omer@example.com\n2,Ali,ali@test.com\n"
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create csv: %v", err)
	}
	defer func() {
		if err := os.Remove(filename); err != nil {
			t.Logf("warning: failed to remove temp file %s: %v", filename, err)
		}
	}()

	ch, err := ReadEmailsFromCSV(filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var results []string
	for email := range ch {
		results = append(results, email)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 emails, got %d", len(results))
	}
}

type errorWriter struct{}

func (e *errorWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("write error")
}

func TestPrintDomainCounts_Error(t *testing.T) {
	writer := &errorWriter{}
	domains := []DomainCount{{Domain: "test.com", Count: 1}}

	err := PrintDomainCounts(domains, writer)
	if err == nil || !strings.Contains(err.Error(), "write error") {
		t.Errorf("expected write error, got %v", err)
	}
}

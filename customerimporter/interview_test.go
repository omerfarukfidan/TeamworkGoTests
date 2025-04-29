package customerimporter

import (
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

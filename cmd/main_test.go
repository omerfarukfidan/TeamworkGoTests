package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestHelpFlag(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--help")
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil && !strings.Contains(output, "Usage of") {
		t.Fatalf("Expected help output, got error: %v\n%s", err, output)
	}
	if !strings.Contains(output, "Usage of") {
		t.Error("Expected help message, not found")
	}
}

func TestMissingInputFile(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	out, _ := cmd.CombinedOutput()
	output := string(out)

	if !strings.Contains(output, "No input file provided") {
		t.Errorf("Expected error for missing input file. Got: %s", output)
	}
}

func TestInvalidFile(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "nonexistent.csv")
	out, _ := cmd.CombinedOutput()
	output := string(out)

	if !strings.Contains(output, "Failed to read CSV file") {
		t.Errorf("Expected error for invalid file. Got: %s", output)
	}
}

func TestValidCSVInputToStdout(t *testing.T) {
	tmpFile := "test_input.csv"
	content := "id,name,email\n1,Omer,omer@example.com\n2,Ali,ali@test.com\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test csv file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			t.Logf("warning: failed to remove temp file %s: %v", tmpFile, err)
		}
	}()

	cmd := exec.Command("go", "run", ".", tmpFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success run, got error: %v\nOutput: %s", err, string(out))
	}
	output := string(out)
	if !strings.Contains(output, "example.com: 1") || !strings.Contains(output, "test.com: 1") {
		t.Errorf("unexpected output: %s", output)
	}
}

func TestValidCSVOutputToFile(t *testing.T) {
	tmpCSV := "test_output.csv"
	tmpOut := "result.txt"
	content := "id,name,email\n1,Omer,omer@example.com\n2,Ali,ali@test.com\n"

	if err := os.WriteFile(tmpCSV, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test csv file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpCSV); err != nil {
			t.Logf("warning: failed to remove temp file %s: %v", tmpCSV, err)
		}
	}()
	defer func() {
		if err := os.Remove(tmpOut); err != nil {
			t.Logf("warning: failed to remove temp file %s: %v", tmpOut, err)
		}
	}()

	cmd := exec.Command("go", "run", ".", "--sort=count", "--output="+tmpOut, tmpCSV)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected no error, got %v\nOutput: %s", err, string(out))
	}

	data, err := os.ReadFile(tmpOut)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(data), "example.com: 1") || !strings.Contains(string(data), "test.com: 1") {
		t.Errorf("unexpected file content: %s", data)
	}
}

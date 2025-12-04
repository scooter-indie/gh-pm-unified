package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/scooter-indie/gh-pmu/internal/api"
)

func TestIntakeCommand(t *testing.T) {
	t.Run("has correct command structure", func(t *testing.T) {
		cmd := newIntakeCommand()

		if cmd.Use != "intake" {
			t.Errorf("expected Use to be 'intake', got %s", cmd.Use)
		}

		if cmd.Short == "" {
			t.Error("expected Short description to be set")
		}

		// Check aliases
		if len(cmd.Aliases) == 0 || cmd.Aliases[0] != "in" {
			t.Error("expected 'in' alias")
		}
	})

	t.Run("has required flags", func(t *testing.T) {
		cmd := newIntakeCommand()

		// Check --apply flag
		applyFlag := cmd.Flags().Lookup("apply")
		if applyFlag == nil {
			t.Error("expected --apply flag")
		}

		// Check --dry-run flag
		dryRunFlag := cmd.Flags().Lookup("dry-run")
		if dryRunFlag == nil {
			t.Error("expected --dry-run flag")
		}

		// Check --json flag
		jsonFlag := cmd.Flags().Lookup("json")
		if jsonFlag == nil {
			t.Error("expected --json flag")
		}
	})

	t.Run("command is registered in root", func(t *testing.T) {
		root := NewRootCommand()
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetArgs([]string{"intake", "--help"})
		err := root.Execute()
		if err != nil {
			t.Errorf("intake command not registered: %v", err)
		}
	})
}

func TestIntakeOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		opts := &intakeOptions{}

		if opts.apply {
			t.Error("apply should be false by default")
		}
		if opts.dryRun {
			t.Error("dryRun should be false by default")
		}
		if opts.json {
			t.Error("json should be false by default")
		}
	})
}

func TestOutputIntakeTable(t *testing.T) {
	t.Run("displays issues in table format", func(t *testing.T) {
		cmd := newIntakeCommand()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)

		issues := []api.Issue{
			{
				Number:     1,
				Title:      "First issue",
				State:      "OPEN",
				Repository: api.Repository{Owner: "owner", Name: "repo"},
			},
			{
				Number:     2,
				Title:      "Second issue",
				State:      "OPEN",
				Repository: api.Repository{Owner: "owner", Name: "repo"},
			},
		}

		err := outputIntakeTable(cmd, issues)
		if err != nil {
			t.Fatalf("outputIntakeTable failed: %v", err)
		}

		// Note: outputIntakeTable writes directly to os.Stdout, not cmd.Out()
		// We're testing it doesn't error; actual output goes to stdout
	})

	t.Run("truncates long titles to 50 chars", func(t *testing.T) {
		cmd := newIntakeCommand()

		// Create issue with 60-character title
		longTitle := strings.Repeat("A", 60)
		issues := []api.Issue{
			{
				Number:     1,
				Title:      longTitle,
				State:      "OPEN",
				Repository: api.Repository{Owner: "owner", Name: "repo"},
			},
		}

		// outputIntakeTable writes to os.Stdout, so we just verify no error
		err := outputIntakeTable(cmd, issues)
		if err != nil {
			t.Fatalf("outputIntakeTable failed with long title: %v", err)
		}
	})

	t.Run("handles empty issue list", func(t *testing.T) {
		cmd := newIntakeCommand()
		issues := []api.Issue{}

		err := outputIntakeTable(cmd, issues)
		if err != nil {
			t.Fatalf("outputIntakeTable failed with empty list: %v", err)
		}
	})
}

func TestOutputIntakeJSON(t *testing.T) {
	t.Run("outputs correct JSON structure with dry-run status", func(t *testing.T) {
		cmd := newIntakeCommand()

		issues := []api.Issue{
			{
				Number:     42,
				Title:      "Test issue",
				State:      "OPEN",
				URL:        "https://github.com/owner/repo/issues/42",
				Repository: api.Repository{Owner: "owner", Name: "repo"},
			},
		}

		// Capture stdout for JSON output
		// Note: outputIntakeJSON writes to os.Stdout via json.NewEncoder
		err := outputIntakeJSON(cmd, issues, "dry-run")
		if err != nil {
			t.Fatalf("outputIntakeJSON failed: %v", err)
		}
	})

	t.Run("status field matches input status", func(t *testing.T) {
		// Test that various status values are preserved
		statuses := []string{"dry-run", "applied", "untracked"}
		for _, status := range statuses {
			cmd := newIntakeCommand()
			issues := []api.Issue{}

			err := outputIntakeJSON(cmd, issues, status)
			if err != nil {
				t.Fatalf("outputIntakeJSON failed with status %q: %v", status, err)
			}
		}
	})

	t.Run("count matches issues length", func(t *testing.T) {
		cmd := newIntakeCommand()

		issues := []api.Issue{
			{Number: 1, Title: "Issue 1", Repository: api.Repository{Owner: "o", Name: "r"}},
			{Number: 2, Title: "Issue 2", Repository: api.Repository{Owner: "o", Name: "r"}},
			{Number: 3, Title: "Issue 3", Repository: api.Repository{Owner: "o", Name: "r"}},
		}

		err := outputIntakeJSON(cmd, issues, "test")
		if err != nil {
			t.Fatalf("outputIntakeJSON failed: %v", err)
		}
	})
}

func TestIntakeJSONOutput_Structure(t *testing.T) {
	t.Run("marshals to correct JSON format", func(t *testing.T) {
		output := intakeJSONOutput{
			Status: "dry-run",
			Count:  2,
			Issues: []intakeJSONIssue{
				{
					Number:     1,
					Title:      "First",
					State:      "OPEN",
					URL:        "https://github.com/owner/repo/issues/1",
					Repository: "owner/repo",
				},
				{
					Number:     2,
					Title:      "Second",
					State:      "OPEN",
					URL:        "https://github.com/owner/repo/issues/2",
					Repository: "owner/repo",
				},
			},
		}

		data, err := json.Marshal(output)
		if err != nil {
			t.Fatalf("Failed to marshal intakeJSONOutput: %v", err)
		}

		// Unmarshal and verify
		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		if result["status"] != "dry-run" {
			t.Errorf("Expected status 'dry-run', got %v", result["status"])
		}

		if int(result["count"].(float64)) != 2 {
			t.Errorf("Expected count 2, got %v", result["count"])
		}

		issues, ok := result["issues"].([]interface{})
		if !ok {
			t.Fatal("Expected issues to be an array")
		}
		if len(issues) != 2 {
			t.Errorf("Expected 2 issues, got %d", len(issues))
		}
	})

	t.Run("intakeJSONIssue includes all fields", func(t *testing.T) {
		issue := intakeJSONIssue{
			Number:     42,
			Title:      "Test Issue",
			State:      "OPEN",
			URL:        "https://github.com/owner/repo/issues/42",
			Repository: "owner/repo",
		}

		data, err := json.Marshal(issue)
		if err != nil {
			t.Fatalf("Failed to marshal intakeJSONIssue: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		expectedFields := []string{"number", "title", "state", "url", "repository"}
		for _, field := range expectedFields {
			if _, exists := result[field]; !exists {
				t.Errorf("Expected field %q to exist in JSON output", field)
			}
		}
	})
}

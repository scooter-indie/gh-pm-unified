package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestCreateCommand_Exists(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"create", "--help"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("create command should exist: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("create")) {
		t.Error("Expected help output to mention 'create'")
	}
}

func TestCreateCommand_HasTitleFlag(t *testing.T) {
	cmd := NewRootCommand()
	createCmd, _, err := cmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}

	flag := createCmd.Flags().Lookup("title")
	if flag == nil {
		t.Fatal("Expected --title flag to exist")
	}
}

func TestCreateCommand_HasBodyFlag(t *testing.T) {
	cmd := NewRootCommand()
	createCmd, _, err := cmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}

	flag := createCmd.Flags().Lookup("body")
	if flag == nil {
		t.Fatal("Expected --body flag to exist")
	}
}

func TestCreateCommand_HasStatusFlag(t *testing.T) {
	cmd := NewRootCommand()
	createCmd, _, err := cmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}

	flag := createCmd.Flags().Lookup("status")
	if flag == nil {
		t.Fatal("Expected --status flag to exist")
	}
}

func TestCreateCommand_HasPriorityFlag(t *testing.T) {
	cmd := NewRootCommand()
	createCmd, _, err := cmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}

	flag := createCmd.Flags().Lookup("priority")
	if flag == nil {
		t.Fatal("Expected --priority flag to exist")
	}
}

func TestCreateCommand_HasLabelFlag(t *testing.T) {
	cmd := NewRootCommand()
	createCmd, _, err := cmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}

	flag := createCmd.Flags().Lookup("label")
	if flag == nil {
		t.Fatal("Expected --label flag to exist")
	}
}

func TestCreateCommand_RequiresTitleInNonInteractiveMode(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"create", "--body", "test body"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when title not provided with --body")
	}
}

// ============================================================================
// createOptions Tests
// ============================================================================

func TestCreateOptions_DefaultValues(t *testing.T) {
	opts := &createOptions{}

	if opts.title != "" {
		t.Errorf("Expected empty title, got %q", opts.title)
	}
	if opts.body != "" {
		t.Errorf("Expected empty body, got %q", opts.body)
	}
	if opts.status != "" {
		t.Errorf("Expected empty status, got %q", opts.status)
	}
	if opts.priority != "" {
		t.Errorf("Expected empty priority, got %q", opts.priority)
	}
	if opts.labels != nil {
		t.Errorf("Expected nil labels, got %v", opts.labels)
	}
}

func TestCreateOptions_WithValues(t *testing.T) {
	opts := &createOptions{
		title:    "Test Issue",
		body:     "Test body content",
		status:   "in_progress",
		priority: "p1",
		labels:   []string{"bug", "urgent"},
	}

	if opts.title != "Test Issue" {
		t.Errorf("Expected title 'Test Issue', got %q", opts.title)
	}
	if opts.body != "Test body content" {
		t.Errorf("Expected body 'Test body content', got %q", opts.body)
	}
	if len(opts.labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(opts.labels))
	}
}

// ============================================================================
// Label Merging Logic Tests
// ============================================================================

func TestLabelMerging_EmptyDefaults(t *testing.T) {
	configLabels := []string{}
	cliLabels := []string{"bug", "urgent"}

	// Simulate the merging logic from runCreate
	labels := append([]string{}, configLabels...)
	labels = append(labels, cliLabels...)

	if len(labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(labels))
	}
	if labels[0] != "bug" || labels[1] != "urgent" {
		t.Errorf("Expected [bug, urgent], got %v", labels)
	}
}

func TestLabelMerging_WithDefaults(t *testing.T) {
	configLabels := []string{"pm-tracked"}
	cliLabels := []string{"bug", "urgent"}

	// Simulate the merging logic from runCreate
	labels := append([]string{}, configLabels...)
	labels = append(labels, cliLabels...)

	if len(labels) != 3 {
		t.Errorf("Expected 3 labels, got %d", len(labels))
	}
	if labels[0] != "pm-tracked" {
		t.Errorf("Expected first label 'pm-tracked', got %q", labels[0])
	}
}

func TestLabelMerging_NoCLILabels(t *testing.T) {
	configLabels := []string{"pm-tracked", "auto-created"}
	var cliLabels []string

	// Simulate the merging logic from runCreate
	labels := append([]string{}, configLabels...)
	labels = append(labels, cliLabels...)

	if len(labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(labels))
	}
}

func TestLabelMerging_BothEmpty(t *testing.T) {
	configLabels := []string{}
	var cliLabels []string

	// Simulate the merging logic from runCreate
	labels := append([]string{}, configLabels...)
	labels = append(labels, cliLabels...)

	if len(labels) != 0 {
		t.Errorf("Expected 0 labels, got %d", len(labels))
	}
}

// ============================================================================
// Error Message Tests
// ============================================================================

func TestCreateCommand_TitleRequiredErrorMessage(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"create"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error when no title provided")
	}

	// The error should mention title is required
	errStr := err.Error()
	if !strings.Contains(errStr, "title") && !strings.Contains(errStr, "configuration") {
		t.Errorf("Expected error about title or config, got: %v", err)
	}
}

func TestCreateCommand_FlagShortcuts(t *testing.T) {
	cmd := NewRootCommand()
	createCmd, _, err := cmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}

	tests := []struct {
		longFlag  string
		shortFlag string
	}{
		{"title", "t"},
		{"body", "b"},
		{"status", "s"},
		{"priority", "p"},
		{"label", "l"},
	}

	for _, tt := range tests {
		t.Run(tt.longFlag, func(t *testing.T) {
			flag := createCmd.Flags().Lookup(tt.longFlag)
			if flag == nil {
				t.Fatalf("Flag --%s not found", tt.longFlag)
			}
			if flag.Shorthand != tt.shortFlag {
				t.Errorf("Expected shorthand -%s for --%s, got -%s", tt.shortFlag, tt.longFlag, flag.Shorthand)
			}
		})
	}
}

func TestCreateCommand_LabelFlagIsArray(t *testing.T) {
	cmd := NewRootCommand()
	createCmd, _, err := cmd.Find([]string{"create"})
	if err != nil {
		t.Fatalf("create command not found: %v", err)
	}

	flag := createCmd.Flags().Lookup("label")
	if flag == nil {
		t.Fatal("Expected --label flag to exist")
	}

	// Check that it's a stringArray type (can be specified multiple times)
	if flag.Value.Type() != "stringArray" {
		t.Errorf("Expected --label to be stringArray, got %s", flag.Value.Type())
	}
}

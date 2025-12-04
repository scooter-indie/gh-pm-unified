package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/scooter-indie/gh-pmu/internal/api"
	"github.com/spf13/cobra"
)

func TestListCommand_Exists(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"list", "--help"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list command should exist: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("list")) {
		t.Error("Expected help output to mention 'list'")
	}
}

func TestListCommand_HasStatusFlag(t *testing.T) {
	cmd := NewRootCommand()
	listCmd, _, err := cmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}

	flag := listCmd.Flags().Lookup("status")
	if flag == nil {
		t.Fatal("Expected --status flag to exist")
	}
}

func TestListCommand_HasPriorityFlag(t *testing.T) {
	cmd := NewRootCommand()
	listCmd, _, err := cmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}

	flag := listCmd.Flags().Lookup("priority")
	if flag == nil {
		t.Fatal("Expected --priority flag to exist")
	}
}

func TestListCommand_HasJSONFlag(t *testing.T) {
	cmd := NewRootCommand()
	listCmd, _, err := cmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}

	flag := listCmd.Flags().Lookup("json")
	if flag == nil {
		t.Fatal("Expected --json flag to exist")
	}
}

func TestListCommand_HasSubIssuesFlag(t *testing.T) {
	cmd := NewRootCommand()
	listCmd, _, err := cmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("list command not found: %v", err)
	}

	flag := listCmd.Flags().Lookup("has-sub-issues")
	if flag == nil {
		t.Fatal("Expected --has-sub-issues flag to exist")
	}

	// Verify it's a boolean flag
	if flag.Value.Type() != "bool" {
		t.Errorf("Expected --has-sub-issues to be bool, got %s", flag.Value.Type())
	}
}

func TestListCommand_HasSubIssuesHelpText(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"list", "--help"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("list help failed: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("has-sub-issues")) {
		t.Error("Expected help to mention --has-sub-issues flag")
	}
}

// ============================================================================
// filterByFieldValue Tests
// ============================================================================

func TestFilterByFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		items     []api.ProjectItem
		fieldName string
		value     string
		wantCount int
	}{
		{
			name: "exact match",
			items: []api.ProjectItem{
				{
					ID: "1",
					FieldValues: []api.FieldValue{
						{Field: "Status", Value: "In Progress"},
					},
				},
				{
					ID: "2",
					FieldValues: []api.FieldValue{
						{Field: "Status", Value: "Done"},
					},
				},
			},
			fieldName: "Status",
			value:     "In Progress",
			wantCount: 1,
		},
		{
			name: "case-insensitive field name",
			items: []api.ProjectItem{
				{
					ID: "1",
					FieldValues: []api.FieldValue{
						{Field: "Status", Value: "Backlog"},
					},
				},
			},
			fieldName: "status",
			value:     "Backlog",
			wantCount: 1,
		},
		{
			name: "case-insensitive value",
			items: []api.ProjectItem{
				{
					ID: "1",
					FieldValues: []api.FieldValue{
						{Field: "Status", Value: "In Progress"},
					},
				},
			},
			fieldName: "Status",
			value:     "in progress",
			wantCount: 1,
		},
		{
			name: "no match",
			items: []api.ProjectItem{
				{
					ID: "1",
					FieldValues: []api.FieldValue{
						{Field: "Status", Value: "Done"},
					},
				},
			},
			fieldName: "Status",
			value:     "In Progress",
			wantCount: 0,
		},
		{
			name:      "empty items",
			items:     []api.ProjectItem{},
			fieldName: "Status",
			value:     "Done",
			wantCount: 0,
		},
		{
			name: "multiple matches",
			items: []api.ProjectItem{
				{
					ID: "1",
					FieldValues: []api.FieldValue{
						{Field: "Priority", Value: "P1"},
					},
				},
				{
					ID: "2",
					FieldValues: []api.FieldValue{
						{Field: "Priority", Value: "P1"},
					},
				},
				{
					ID: "3",
					FieldValues: []api.FieldValue{
						{Field: "Priority", Value: "P2"},
					},
				},
			},
			fieldName: "Priority",
			value:     "P1",
			wantCount: 2,
		},
		{
			name: "item with multiple fields",
			items: []api.ProjectItem{
				{
					ID: "1",
					FieldValues: []api.FieldValue{
						{Field: "Status", Value: "Done"},
						{Field: "Priority", Value: "P1"},
					},
				},
			},
			fieldName: "Priority",
			value:     "P1",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByFieldValue(tt.items, tt.fieldName, tt.value)
			if len(result) != tt.wantCount {
				t.Errorf("filterByFieldValue() returned %d items, want %d", len(result), tt.wantCount)
			}
		})
	}
}

// ============================================================================
// getFieldValue Tests
// ============================================================================

func TestGetFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		item      api.ProjectItem
		fieldName string
		want      string
	}{
		{
			name: "field exists",
			item: api.ProjectItem{
				FieldValues: []api.FieldValue{
					{Field: "Status", Value: "In Progress"},
				},
			},
			fieldName: "Status",
			want:      "In Progress",
		},
		{
			name: "field missing",
			item: api.ProjectItem{
				FieldValues: []api.FieldValue{
					{Field: "Status", Value: "Done"},
				},
			},
			fieldName: "Priority",
			want:      "",
		},
		{
			name: "case-insensitive lookup",
			item: api.ProjectItem{
				FieldValues: []api.FieldValue{
					{Field: "Priority", Value: "P0"},
				},
			},
			fieldName: "priority",
			want:      "P0",
		},
		{
			name: "multiple fields returns first match",
			item: api.ProjectItem{
				FieldValues: []api.FieldValue{
					{Field: "Status", Value: "Done"},
					{Field: "Priority", Value: "P1"},
					{Field: "Size", Value: "M"},
				},
			},
			fieldName: "Size",
			want:      "M",
		},
		{
			name:      "empty field values",
			item:      api.ProjectItem{},
			fieldName: "Status",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFieldValue(tt.item, tt.fieldName)
			if result != tt.want {
				t.Errorf("getFieldValue() = %q, want %q", result, tt.want)
			}
		})
	}
}

// ============================================================================
// outputTable Tests
// ============================================================================

// createTestCmd creates a cobra command with output set to a buffer
func createTestCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	return cmd
}

func TestOutputTable_EmptyItems(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := createTestCmd(buf)

	err := outputTable(cmd, []api.ProjectItem{})
	if err != nil {
		t.Fatalf("outputTable() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No issues found") {
		t.Errorf("Expected 'No issues found', got: %s", output)
	}
}

func TestOutputTable_TitleTruncation(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := createTestCmd(buf)

	longTitle := "This is a very long title that exceeds fifty characters and should be truncated"
	items := []api.ProjectItem{
		{
			ID: "1",
			Issue: &api.Issue{
				Number: 1,
				Title:  longTitle,
				State:  "OPEN",
			},
			FieldValues: []api.FieldValue{
				{Field: "Status", Value: "Done"},
			},
		},
	}

	// Note: outputTable writes to os.Stdout, not cmd.Out()
	// We can't capture this directly, but we can verify no error
	err := outputTable(cmd, items)
	if err != nil {
		t.Fatalf("outputTable() error = %v", err)
	}
}

func TestOutputTable_WithAssignees(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := createTestCmd(buf)

	items := []api.ProjectItem{
		{
			ID: "1",
			Issue: &api.Issue{
				Number: 42,
				Title:  "Test Issue",
				State:  "OPEN",
				Assignees: []api.Actor{
					{Login: "user1"},
					{Login: "user2"},
				},
			},
			FieldValues: []api.FieldValue{
				{Field: "Status", Value: "In Progress"},
				{Field: "Priority", Value: "P1"},
			},
		},
	}

	err := outputTable(cmd, items)
	if err != nil {
		t.Fatalf("outputTable() error = %v", err)
	}
}

func TestOutputTable_NoAssignees(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := createTestCmd(buf)

	items := []api.ProjectItem{
		{
			ID: "1",
			Issue: &api.Issue{
				Number: 1,
				Title:  "No Assignee Issue",
				State:  "OPEN",
			},
		},
	}

	err := outputTable(cmd, items)
	if err != nil {
		t.Fatalf("outputTable() error = %v", err)
	}
}

func TestOutputTable_NilIssue(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := createTestCmd(buf)

	items := []api.ProjectItem{
		{ID: "1", Issue: nil},
		{
			ID: "2",
			Issue: &api.Issue{
				Number: 1,
				Title:  "Valid Issue",
				State:  "OPEN",
			},
		},
	}

	err := outputTable(cmd, items)
	if err != nil {
		t.Fatalf("outputTable() error = %v", err)
	}
}

// ============================================================================
// outputJSON Tests
// ============================================================================

func TestOutputJSON_EmptyItems(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := createTestCmd(buf)

	// outputJSON writes to os.Stdout, not cmd buffer
	// But we can verify structure by checking for error
	err := outputJSON(cmd, []api.ProjectItem{})
	if err != nil {
		t.Fatalf("outputJSON() error = %v", err)
	}
}

func TestOutputJSON_WithItems(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := createTestCmd(buf)

	items := []api.ProjectItem{
		{
			ID: "1",
			Issue: &api.Issue{
				Number: 42,
				Title:  "Test Issue",
				State:  "OPEN",
				URL:    "https://github.com/owner/repo/issues/42",
				Repository: api.Repository{
					Owner: "owner",
					Name:  "repo",
				},
				Assignees: []api.Actor{
					{Login: "user1"},
				},
			},
			FieldValues: []api.FieldValue{
				{Field: "Status", Value: "In Progress"},
				{Field: "Priority", Value: "P1"},
			},
		},
	}

	err := outputJSON(cmd, items)
	if err != nil {
		t.Fatalf("outputJSON() error = %v", err)
	}
}

func TestOutputJSON_NilIssue(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := createTestCmd(buf)

	items := []api.ProjectItem{
		{ID: "1", Issue: nil},
	}

	err := outputJSON(cmd, items)
	if err != nil {
		t.Fatalf("outputJSON() error = %v", err)
	}
}

func TestJSONOutput_Structure(t *testing.T) {
	// Test that JSONOutput struct has expected fields
	output := JSONOutput{
		Items: []JSONItem{
			{
				Number:     1,
				Title:      "Test",
				State:      "OPEN",
				URL:        "https://example.com",
				Repository: "owner/repo",
				Assignees:  []string{"user1"},
				FieldValues: map[string]string{
					"Status": "Done",
				},
			},
		},
	}

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("Failed to marshal JSONOutput: %v", err)
	}

	// Verify it can be unmarshaled back
	var parsed JSONOutput
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSONOutput: %v", err)
	}

	if len(parsed.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(parsed.Items))
	}
	if parsed.Items[0].Number != 1 {
		t.Errorf("Expected number 1, got %d", parsed.Items[0].Number)
	}
	if parsed.Items[0].FieldValues["Status"] != "Done" {
		t.Errorf("Expected Status=Done, got %s", parsed.Items[0].FieldValues["Status"])
	}
}

func TestJSONItem_AllFields(t *testing.T) {
	item := JSONItem{
		Number:      42,
		Title:       "Test Issue",
		State:       "OPEN",
		URL:         "https://github.com/owner/repo/issues/42",
		Repository:  "owner/repo",
		Assignees:   []string{"user1", "user2"},
		FieldValues: map[string]string{"Status": "Done", "Priority": "P1"},
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal JSONItem: %v", err)
	}

	jsonStr := string(data)
	expectedFields := []string{"number", "title", "state", "url", "repository", "assignees", "fieldValues"}
	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("Expected JSON to contain field %q", field)
		}
	}
}

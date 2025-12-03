# gh-pm-unified

A unified GitHub CLI extension for project management and sub-issue hierarchy.

## Features

- **Project Management**: List, view, create, and update issues with project field values
- **Sub-Issue Hierarchy**: Create and manage parent-child issue relationships
- **Batch Operations**: Intake untracked issues and triage with configurable rules
- **Issue Splitting**: Convert checklists into sub-issues automatically

## Installation

### From GitHub Releases

```bash
gh extension install scooter-indie/gh-pm-unified
```

After installation, use the extension as `gh pm-unified`:

```bash
gh pm-unified --help
```

### From Source

```bash
git clone https://github.com/scooter-indie/gh-pm-unified.git
cd gh-pm-unified
go build -o gh-pm-unified .
gh extension install .
```

## Quick Start

1. Initialize configuration in your repository:

```bash
gh pm-unified init
```

2. List issues with project metadata:

```bash
gh pm-unified list
```

3. View an issue with all project fields:

```bash
gh pm-unified view 123
```

## Commands

```
gh pm-unified [command]

Project Management:
  init        Initialize configuration
  list        List issues with project metadata
  view        View issue with project fields
  create      Create issue with project fields
  move        Update issue project fields

Sub-Issue Management:
  sub add     Link existing issue as sub-issue
  sub create  Create new sub-issue under parent
  sub list    List sub-issues of a parent
  sub remove  Unlink sub-issue from parent

Batch Operations:
  intake      Find and add untracked issues to project
  triage      Bulk update issues based on config rules
  split       Create sub-issues from checklist or arguments

Flags:
  -h, --help      help for gh-pm-unified
  -v, --version   version for gh-pm-unified
```

## Configuration

gh-pm-unified uses a `.gh-pm.yml` file in your repository root:

```yaml
project:
  owner: your-username
  number: 1

repositories:
  - your-username/your-repo

defaults:
  priority: P2
  status: Backlog

# Triage rules for batch operations
triage:
  stale-issues:
    query: "is:open updated:<2024-01-01"
    apply:
      labels:
        - needs-triage
```

## Command Examples

### Project Management

```bash
# Initialize project configuration interactively
gh pm-unified init

# List all issues in project
gh pm-unified list

# List issues filtered by status
gh pm-unified list --status "In Progress"

# View issue with project fields
gh pm-unified view 42

# Create issue with project fields
gh pm-unified create --title "New feature" --status "Backlog" --priority "P1"

# Update issue status
gh pm-unified move 42 --status "In Progress"
```

### Sub-Issue Management

```bash
# Add existing issue as sub-issue
gh pm-unified sub add 10 15  # Issue 15 becomes sub-issue of 10

# Create new sub-issue
gh pm-unified sub create --parent 10 --title "Subtask 1"

# List sub-issues
gh pm-unified sub list 10

# Remove sub-issue link
gh pm-unified sub remove 10 15
```

### Batch Operations

```bash
# Find untracked issues
gh pm-unified intake --dry-run

# Add untracked issues to project
gh pm-unified intake --apply

# Run triage rule
gh pm-unified triage stale-issues --dry-run

# Split issue from checklist in body
gh pm-unified split 42 --from body

# Split issue from arguments
gh pm-unified split 42 "Task 1" "Task 2" "Task 3"
```

## Development

### Prerequisites

- Go 1.21+
- GitHub CLI (`gh`) with `project` scope

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Integration Tests

```bash
# Requires authenticated gh CLI
go test -tags=integration ./internal/api/...
```

## Project Status

### Completed

**Epic 1: Core Unification** ✅
- Project configuration initialization
- Issue listing with project metadata
- Issue viewing with project fields
- Issue creation with project fields
- Issue field updates (move)
- Sub-issue management (add, create, list, remove)
- Issue intake (find untracked)
- Issue triage (batch updates)
- Issue splitting

**Epic 3: Enhanced Integration** ✅
- Native sub-issue handling in split
- Cross-repository sub-issues (`gh pm-unified sub create --repo owner/repo`)
- Sub-issue progress tracking (`gh pm-unified view` shows progress bar)
- Recursive operations (`gh pm-unified move --recursive`)

### Not Implemented
- **Epic 2: Project Templates** - GitHub API does not support creating views programmatically
- **Epic 4: Template Ecosystem** - Blocked by Epic 2

For project management operations like copying projects, use native `gh project` commands:
```bash
gh project copy 17 --source-owner myorg --target-owner myorg --title "New Project"
```

## Attribution

This project builds upon and unifies functionality from the following open-source projects by [@yahsan2](https://github.com/yahsan2):

- **[gh-pm](https://github.com/yahsan2/gh-pm)** - GitHub project management CLI extension
- **[gh-sub-issue](https://github.com/yahsan2/gh-sub-issue)** - Sub-issue hierarchy management

Thank you to the original author for the foundational work.

## License

MIT

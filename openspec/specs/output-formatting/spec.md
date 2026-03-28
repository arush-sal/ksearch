# Output Formatting Specification

## Purpose
Manages rendering of Kubernetes resource lists into human-readable
tab-aligned table output on stdout.

## Requirements

### Requirement: Printer output target
Each printer function SHALL accept an io.Writer and write to it,
allowing callers to capture output into a buffer before flushing.

#### Scenario: Output captured into buffer
- GIVEN a Printer call with a bytes.Buffer as the writer
- WHEN the resource is rendered
- THEN the formatted table is written into the buffer, not directly to stdout

### Requirement: Pattern filter ownership
A single matchesPattern(name, pattern string) bool helper SHALL own
all name-pattern filtering; printer functions SHALL delegate to it.

#### Scenario: Pattern filter applied uniformly
- GIVEN --pattern=nginx is passed
- WHEN any resource type is printed
- THEN only items whose name contains "nginx" are shown
- AND the filter logic exists in exactly one place in the code

### Requirement: Deterministic output order
The system SHALL output resource sections in the same order as the
resource kinds list, regardless of the order in which goroutines complete.

#### Scenario: Ordered output under concurrency
- GIVEN Pods and Deployments are rendered concurrently
- WHEN both goroutines finish
- THEN Pods section always appears before Deployments section in stdout

### Requirement: Non-empty sections only
The system SHALL omit a resource section entirely if no items matched.

#### Scenario: Empty section suppressed
- GIVEN a namespace with no Secrets
- WHEN ksearch is invoked
- THEN no "Secrets" header or table is printed

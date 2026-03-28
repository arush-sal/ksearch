# Proposal: Concurrent Printing and Printers Refactor

## Purpose
Currently results are fetched concurrently via a goroutine but printed
serially in the main loop, one resource type at a time. `printers.go`
contains 17 nearly-identical functions each duplicating the same
pattern-match guard. This change parallelises the print phase and
eliminates the duplication through a cleaner, table-driven printer design.

## Requirements

### Requirement: Concurrent output assembly
The tool SHALL assemble output for all resource types concurrently.

#### Scenario: Multiple resource types rendered in parallel
- GIVEN a namespace with Pods, Deployments, and ConfigMaps
- WHEN ksearch is invoked
- THEN all three resource sections are rendered concurrently (not sequentially)
- AND the final stdout output is deterministically ordered by resource type

### Requirement: No interleaved output
Concurrent rendering SHALL NOT produce interleaved or garbled output on stdout.

#### Scenario: Clean stdout under concurrency
- GIVEN two resource types are rendered simultaneously by separate goroutines
- WHEN both goroutines complete
- THEN stdout contains each resource section exactly once, cleanly separated

### Requirement: Printers refactor — no duplicated filter logic
Each printer function SHALL NOT implement its own name-pattern filter.
A single shared helper SHALL own the filter check.

#### Scenario: Pattern filter applied uniformly
- GIVEN --pattern=nginx is passed
- WHEN any resource type is printed
- THEN only items whose name contains "nginx" are shown
- AND the strings.Contains check exists in exactly one place in the codebase

### Requirement: Printer writes to a caller-supplied writer
Each printer function SHALL accept an io.Writer so that callers can
buffer output before flushing to stdout.

#### Scenario: Printer writes to a buffer
- GIVEN a bytes.Buffer is passed as the writer
- WHEN Printer() is called for a ConfigMapList
- THEN the formatted table is written into the buffer, not stdout

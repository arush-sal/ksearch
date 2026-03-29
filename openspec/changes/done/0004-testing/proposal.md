# Proposal: Test Suite

## Purpose
The repository currently has zero test files. This change establishes the
full test suite covering `pkg/printers`, `pkg/util`, and `pkg/cache`,
including a mandatory security regression test that asserts secret values
never appear in printer output or cache files.

## Requirements

### Requirement: Printer security regression test
The test suite SHALL include a test that passes known secret values through
the Printer and asserts they are absent from output.

#### Scenario: Secret values never printed
- GIVEN a SecretList with Data containing sentinel strings
- WHEN Printer() is called
- THEN the sentinel strings do not appear in output
- AND the data count is correctly shown

### Requirement: Pattern filter unit tests
The test suite SHALL verify that the pattern filter includes matching
items and excludes non-matching items across resource types.

#### Scenario: Filter excludes non-matching items
- GIVEN a list with items "nginx-config" and "redis-config"
- WHEN Printer() is called with pattern="nginx"
- THEN only "nginx-config" appears in output

### Requirement: Getter channel contract tests
The test suite SHALL verify that Getter() always closes its output
channel, including on unknown kinds and API errors.

#### Scenario: Channel closed on unknown kind
- GIVEN an unrecognised kind string passed to Getter()
- WHEN the goroutine completes
- THEN the channel is closed (no deadlock)

### Requirement: Cache security test
The test suite SHALL include a test that writes a SectionEntry to disk
and scans the raw JSON file for known sensitive strings.

#### Scenario: Cache file contains no sensitive data
- GIVEN a SectionEntry with printer-formatted output
- WHEN Write() is called
- THEN the raw cache file contains no sentinel strings
